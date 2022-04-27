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
	state        variableState
	name         Token
	nameType     nameType
	functionDecl *StmtFuncDecl
}

type checker struct {
	lines  [][]rune
	scopes []map[string]variable
	scope  int
	state  map[string]any
}

func (c *checker) copyState() map[string]any {
	oldState := c.state
	c.state = make(map[string]any, len(c.state))
	for k, v := range oldState {
		c.state[k] = v
	}
	return oldState
}

func Check(program []Stmt, lines [][]rune) error {
	checker := &checker{
		lines:  lines,
		scopes: make([]map[string]variable, 0),
		scope:  -1,
	}
	checker.beginScope()

	checker.state = map[string]any{
		"inLoop":           false,
		"returnValueCount": 0,
		"canThrow":         false,
		"inTry":            false,
	}

	for name, callable := range nativeFunctions {
		checker.scopes[checker.scope][name] = variable{
			state:    variableStateUsed,
			nameType: nameTypeFunction,
			functionDecl: &StmtFuncDecl{
				Name: Token{
					Lexeme: name,
					Line:   -1,
					Column: -1,
					Type:   IDENTIFIER,
				},
				ReturnValueCount: callable.ArgumentCount(),
				Throws:           callable.Throws(),
			},
		}
	}

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
		name:         stmt.Name,
		state:        state,
		nameType:     nameTypeFunction,
		functionDecl: stmt,
	}

	c.beginScope()
	defer c.endScope()
	for _, p := range stmt.Parameters {
		c.scopes[c.scope][p] = variable{
			state:    variableStateUsed,
			nameType: nameTypeVariable,
		}
	}

	oldState := c.copyState()
	c.state["inLoop"] = false
	c.state["returnValueCount"] = stmt.ReturnValueCount
	c.state["canThrow"] = stmt.Throws

	err := stmt.Body.Accept(c)

	c.state = oldState

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

	oldState := c.copyState()
	c.state["inLoop"] = true
	err = stmt.Body.Accept(c)
	c.state = oldState
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

	oldState := c.copyState()
	c.state["inLoop"] = true
	err = stmt.Body.Accept(c)
	c.state = oldState
	if err != nil {
		return err
	}

	return nil
}

func (c *checker) VisitLoopControl(stmt *StmtLoopControl) error {
	if !c.state["inLoop"].(bool) {
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
	if len(stmt.Values) != c.state["returnValueCount"].(int) {
		return c.newError(fmt.Sprintf("Wrong return value count. Expected %d, got %d.", c.state["returnValueCount"].(int), len(stmt.Values)), stmt.Keyword)
	}
	for _, v := range stmt.Values {
		_, err := v.Accept(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *checker) VisitThrow(stmt *StmtThrow) error {
	if !c.state["canThrow"].(bool) {
		return c.newError("Cannot throw exception in non-throwing function. Append 'throws' to the function signature.", stmt.Keyword)
	}
	return nil
}

func (c *checker) VisitTry(stmt *StmtTry) error {
	oldState := c.copyState()
	c.state["inTry"] = true
	err := stmt.Body.Accept(c)
	c.state = oldState
	if err != nil {
		return err
	}

	c.beginScope()
	defer c.endScope()

	if stmt.ExceptionName.Lexeme != "" {
		c.scopes[c.scope][stmt.ExceptionName.Lexeme] = variable{
			name:     stmt.ExceptionName,
			state:    variableStateDeclared,
			nameType: nameTypeVariable,
		}
	}

	err = stmt.CatchBody.Accept(c)
	if err != nil {
		return err
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
	var returnValueCount any
	if v, ok := expr.Callee.(*ExprVariable); ok {
		scope := c.findVariable(v.Name.Lexeme)
		if scope < 0 {
			return nil, c.newError("Undefined name.", v.Name)
		}
		variable := c.scopes[scope][v.Name.Lexeme]

		if variable.nameType == nameTypeFunction && variable.functionDecl != nil {
			if variable.functionDecl.Throws && !c.state["canThrow"].(bool) && !c.state["inTry"].(bool) {
				return nil, c.newError("Calling throwing function in a non-throwing function outside of a try block.", v.Name)
			}
			returnValueCount = variable.functionDecl.ReturnValueCount
		}
	}

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

	return returnValueCount, nil
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
		ret, err := assignee.Accept(c)
		if err != nil {
			return nil, err
		}
		if returnValueCount, ok := ret.(int); ok {
			if returnValueCount != len(assign.Assignees) {
				return nil, c.newError(fmt.Sprintf("Cannot assign %d values to %d variables.", returnValueCount, assign.Assignees), assign.Operator)
			}
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
