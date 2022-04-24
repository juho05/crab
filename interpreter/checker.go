package interpreter

import "fmt"

type variableState int

const (
	variableStateDeclared variableState = iota
	variableStateDefined
	variableStateUsed
)

type nameType string

const (
	nameTypeVariable nameType = "variable"
	nameTypeFunction nameType = "function"
)

type variable struct {
	state    variableState
	name     Token
	nameType nameType
}

type checker struct {
	lines            [][]rune
	scopes           []map[string]variable
	scope            int
	inLoop           bool
	returnValueCount int
}

func Check(program []Stmt, lines [][]rune) error {
	checker := &checker{
		lines:  lines,
		scopes: make([]map[string]variable, 0),
		scope:  -1,
	}
	checker.beginScope()
	registerNativeFunctions(checker.scopes[checker.scope])

	for _, stmt := range program {
		err := stmt.Accept(checker)
		if err != nil {
			return err
		}
	}

	for _, scope := range checker.scopes {
		for _, v := range scope {
			if v.state != variableStateUsed {
				fmt.Println(generateWarningText(fmt.Sprintf("Unused %s.", v.nameType), lines[v.name.Line], v.name.Line, v.name.Column, v.name.Column+len([]byte(v.name.Lexeme))))
			}
		}
	}

	return nil
}

func (c *checker) VisitExpression(stmt *StmtExpression) error {
	_, err := stmt.Expr.Accept(c)
	return err
}

func (c *checker) VisitVarDecl(stmt *StmtVarDecl) error {
	for _, name := range stmt.Names {
		if _, ok := c.scopes[c.scope][name.Lexeme]; ok {
			return c.newError(fmt.Sprintf("'%s' is already defined in this scope", name.Lexeme), name)
		}
		c.scopes[c.scope][name.Lexeme] = variable{
			name:     name,
			state:    variableStateDeclared,
			nameType: nameTypeVariable,
		}
	}
	_, err := stmt.Expr.Accept(c)
	if err != nil {
		return err
	}
	for _, name := range stmt.Names {
		c.scopes[c.scope][name.Lexeme] = variable{
			name:     name,
			state:    variableStateDefined,
			nameType: nameTypeVariable,
		}
	}

	return nil
}

func (c *checker) VisitFuncDecl(stmt *StmtFuncDecl) error {
	if _, ok := c.scopes[c.scope][stmt.Name.Lexeme]; ok {
		return c.newError(fmt.Sprintf("'%s' is already defined in this scope", stmt.Name.Lexeme), stmt.Name)
	}

	state := variableStateDefined
	if c.scope == 0 && stmt.Name.Lexeme == "main" {
		state = variableStateUsed
	}

	c.scopes[c.scope][stmt.Name.Lexeme] = variable{
		name:     stmt.Name,
		state:    state,
		nameType: nameTypeFunction,
	}

	c.beginScope()
	defer c.endScope()
	for _, p := range stmt.Parameters {
		c.scopes[c.scope][p] = variable{
			state:    variableStateUsed,
			nameType: nameTypeVariable,
		}
	}

	wasInLoop := c.inLoop
	c.inLoop = false

	prevReturnValueCount := c.returnValueCount
	c.returnValueCount = stmt.ReturnValueCount

	err := stmt.Body.Accept(c)

	c.inLoop = wasInLoop
	c.returnValueCount = prevReturnValueCount

	return err
}

func (c *checker) VisitIf(stmt *StmtIf) error {
	_, err := stmt.Condition.Accept(c)
	if err != nil {
		return err
	}

	err = stmt.Body.Accept(c)
	if err != nil {
		return err
	}

	if stmt.ElseBody != nil {
		err = stmt.ElseBody.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) VisitWhile(stmt *StmtWhile) error {
	_, err := stmt.Condition.Accept(c)
	if err != nil {
		return err
	}

	wasInLoop := c.inLoop
	c.inLoop = true
	err = stmt.Body.Accept(c)
	c.inLoop = wasInLoop
	if err != nil {
		return err
	}

	return nil
}

func (c *checker) VisitFor(stmt *StmtFor) error {
	err := stmt.Initializer.Accept(c)
	if err != nil {
		return err
	}
	_, err = stmt.Increment.Accept(c)
	if err != nil {
		return err
	}
	_, err = stmt.Condition.Accept(c)
	if err != nil {
		return err
	}

	wasInLoop := c.inLoop
	c.inLoop = true
	err = stmt.Body.Accept(c)
	c.inLoop = wasInLoop
	if err != nil {
		return err
	}

	return nil
}

func (c *checker) VisitLoopControl(stmt *StmtLoopControl) error {
	if !c.inLoop {
		switch stmt.Keyword.Type {
		case BREAK:
			return c.newError("'break' statement outside of loop.", stmt.Keyword)
		case CONTINUE:
			return c.newError("'continue' statement outside of loop.", stmt.Keyword)
		}
	}
	return nil
}

func (c *checker) VisitReturn(stmt *StmtReturn) error {
	if len(stmt.Values) != c.returnValueCount {
		return c.newError(fmt.Sprintf("Wrong return value count. Expected %d, got %d.", c.returnValueCount, len(stmt.Values)), stmt.Keyword)
	}
	for _, v := range stmt.Values {
		_, err := v.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) VisitBlock(stmt *StmtBlock) error {
	c.beginScope()
	defer c.endScope()
	for _, s := range stmt.Statements {
		err := s.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) VisitLiteral(expr *ExprLiteral) (any, error) {
	return nil, nil
}

func (c *checker) VisitVariable(expr *ExprVariable) (any, error) {
	scope := c.findVariable(expr.Name.Lexeme)
	if scope < 0 {
		return nil, c.newError("Undefined name.", expr.Name)
	}
	expr.NestingLevel = scope

	v := c.scopes[scope][expr.Name.Lexeme]
	c.scopes[scope][expr.Name.Lexeme] = variable{
		name:     v.name,
		nameType: v.nameType,
		state:    variableStateUsed,
	}

	return nil, nil
}

func (c *checker) VisitCall(expr *ExprCall) (any, error) {
	_, err := expr.Callee.Accept(c)
	if err != nil {
		return nil, err
	}

	for _, a := range expr.Args {
		_, err = a.Accept(c)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (c *checker) VisitSubscript(expr *ExprSubscript) (any, error) {
	_, err := expr.Object.Accept(c)
	if err != nil {
		return nil, err
	}
	return expr.Subscript.Accept(c)
}

func (c *checker) VisitGrouping(expr *ExprGrouping) (any, error) {
	return expr.Expr.Accept(c)
}

func (c *checker) VisitList(list *ExprList) (any, error) {
	for _, v := range list.Values {
		_, err := v.Accept(c)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (c *checker) VisitUnary(expr *ExprUnary) (any, error) {
	return expr.Right.Accept(c)
}

func (c *checker) VisitBinary(expr *ExprBinary) (any, error) {
	_, err := expr.Left.Accept(c)
	if err != nil {
		return nil, err
	}
	return expr.Right.Accept(c)
}

func (c *checker) VisitLogical(expr *ExprLogical) (any, error) {
	_, err := expr.Left.Accept(c)
	if err != nil {
		return nil, err
	}
	return expr.Right.Accept(c)
}

func (c *checker) VisitTernary(expr *ExprTernary) (any, error) {
	_, err := expr.Left.Accept(c)
	if err != nil {
		return nil, err
	}
	_, err = expr.Center.Accept(c)
	if err != nil {
		return nil, err
	}
	return expr.Right.Accept(c)
}

func (c *checker) VisitAssign(assign *ExprAssign) (any, error) {
	for _, assignee := range assign.Assignees {
		_, err := assignee.Accept(c)
		if err != nil {
			return nil, err
		}
	}
	return assign.Expr.Accept(c)
}

func (c *checker) beginScope() {
	c.scopes = append(c.scopes, make(map[string]variable))
	c.scope++
}

func (c *checker) endScope() {
	c.scope--
	c.scopes = c.scopes[:len(c.scopes)-1]
}

func (c *checker) findVariable(name string) int {
	scope := c.scope
	for scope >= 0 {
		if _, ok := c.scopes[scope][name]; ok {
			break
		}
		scope--
	}
	return scope
}

func (c *checker) newError(message string, token Token) error {
	return ParseError{
		Token:   token,
		Message: message,
		Line:    c.lines[token.Line],
	}
}
