package interpreter

import "fmt"

type variableState int

const (
	variableStateDeclared variableState = iota
	variableStateDefined
	variableStateUsed
)

type checker struct {
	lines  [][]rune
	scopes []map[string]variableState
	scope  int
}

func Check(program []Stmt, lines [][]rune) error {
	checker := &checker{
		lines:  lines,
		scopes: make([]map[string]variableState, 0),
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

	return nil
}

func (c *checker) VisitExpression(stmt *StmtExpression) error {
	_, err := stmt.Expr.Accept(c)
	return err
}

func (c *checker) VisitVarDecl(stmt *StmtVarDecl) error {
	if _, ok := c.scopes[c.scope][stmt.Name.Lexeme]; ok {
		return c.newError(fmt.Sprintf("'%s' is already defined in this scope", stmt.Name.Lexeme), stmt.Name)
	}

	c.scopes[c.scope][stmt.Name.Lexeme] = variableStateDeclared
	_, err := stmt.Expr.Accept(c)
	if err != nil {
		return err
	}
	c.scopes[c.scope][stmt.Name.Lexeme] = variableStateDefined
	return nil
}

func (c *checker) VisitFuncDecl(stmt *StmtFuncDecl) error {
	if _, ok := c.scopes[c.scope][stmt.Name.Lexeme]; ok {
		return c.newError(fmt.Sprintf("'%s' is already defined in this scope", stmt.Name.Lexeme), stmt.Name)
	}

	c.scopes[c.scope][stmt.Name.Lexeme] = variableStateDeclared
	c.scopes[c.scope][stmt.Name.Lexeme] = variableStateDefined

	c.beginScope()
	defer c.endScope()
	for _, p := range stmt.Parameters {
		c.scopes[c.scope][p] = variableStateUsed
	}

	return stmt.Body.Accept(c)
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

func (c *checker) VisitGrouping(expr *ExprGrouping) (any, error) {
	return expr.Expr.Accept(c)
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

func (c *checker) beginScope() {
	c.scopes = append(c.scopes, make(map[string]variableState))
	c.scope++
}

func (c *checker) endScope() {
	c.scope--
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
