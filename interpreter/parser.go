package interpreter

import (
	"fmt"
)

type parser struct {
	tokens  []Token
	current int
	lines   [][]rune
	errors  []error
}

func Parse(tokens []Token, lines [][]rune) ([]Stmt, []error) {
	parser := &parser{
		tokens: tokens,
		lines:  lines,
		errors: make([]error, 0),
	}
	return parser.parse()
}

func (p *parser) parse() ([]Stmt, []error) {
	statements := make([]Stmt, 0)
	for p.peek().Type != EOF {
		statements = append(statements, p.declaration(false))
	}
	return statements, p.errors
}

func (p *parser) declaration(allowNonDeclarationStatements bool) Stmt {
	var stmt Stmt
	var err error
	if p.match(VAR) {
		stmt, err = p.varDecl()
	} else if p.match(FUNC) {
		stmt, err = p.funcDecl()
	} else if allowNonDeclarationStatements {
		stmt, err = p.statement()
	}

	if stmt == nil && err == nil {
		err = p.newError(fmt.Sprintf("Unexpected token '%s'", p.peek().Lexeme))
	}

	if err != nil {
		p.errors = append(p.errors, err)
		p.synchronize()
	}
	return stmt
}

func (p *parser) varDecl() (Stmt, error) {
	if !p.match(IDENTIFIER) {
		return nil, p.newError("Expect identifier after 'var' keyword.")
	}

	name := p.previous()

	var expr Expr
	var err error
	if p.match(EQUAL) {
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.match(SEMICOLON) {
		return nil, p.newError("Missing semicolon.")
	}

	return &StmtVarDecl{
		Name: name,
		Expr: expr,
	}, nil
}

func (p *parser) funcDecl() (Stmt, error) {
	if !p.match(IDENTIFIER) {
		return nil, p.newError("Expect identifier after 'func' keyword.")
	}

	name := p.previous()

	if !p.match(OPEN_PAREN) {
		return nil, p.newError("Expect '(' after function name.")
	}

	parameters := make([]string, 0)
	for p.peek().Type != CLOSE_PAREN {
		if !p.match(IDENTIFIER) {
			return nil, p.newError("Invalid parameter name.")
		}
		parameters = append(parameters, p.previous().Lexeme)
		if p.peek().Type == CLOSE_PAREN {
			break
		}
		if !p.match(COMMA) {
			return nil, p.newError("Expect ',' between parameters.")
		}
	}

	if !p.match(CLOSE_PAREN) {
		return nil, p.newError("Expect ')' after function parameter list.")
	}

	if !p.match(OPEN_BRACE) {
		return nil, p.newError("Expect block after function signature.")
	}

	block, err := p.block()
	if err != nil {
		return nil, err
	}

	return &StmtFuncDecl{
		Name:       name,
		Body:       block,
		Parameters: parameters,
	}, nil
}

func (p *parser) statement() (Stmt, error) {
	if p.match(OPEN_BRACE) {
		return p.block()
	}
	if p.match(IF) {
		return p.ifStmt()
	}
	return p.expressionStmt()
}

func (p *parser) ifStmt() (Stmt, error) {
	if !p.match(OPEN_PAREN) {
		return nil, p.newError("Expect '(' after 'if'.")
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(CLOSE_PAREN) {
		return nil, p.newError("Expect ')' after if condition.")
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBody Stmt
	if p.match(ELSE) {
		elseBody, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &StmtIf{
		Condition: condition,
		Body:      body,
		ElseBody:  elseBody,
	}, nil
}

func (p *parser) expressionStmt() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(SEMICOLON) {
		return nil, p.newError("Missing semicolon.")
	}

	return &StmtExpression{
		Expr: expr,
	}, nil
}

func (p *parser) block() (Stmt, error) {
	openBrace := p.previous()
	statements := make([]Stmt, 0)

	for p.peek().Type != CLOSE_BRACE && p.peek().Type != EOF {
		stmt := p.declaration(true)
		statements = append(statements, stmt)
	}

	if !p.match(CLOSE_BRACE) {
		return nil, p.newErrorAt("Block never closed.", openBrace)
	}

	return &StmtBlock{
		Statements: statements,
	}, nil
}

func (p *parser) expression() (Expr, error) {
	return p.assign()
}

func (p *parser) assign() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		if v, ok := expr.(*ExprVariable); ok {
			right, err := p.assign()
			if err != nil {
				return nil, err
			}
			return &ExprAssign{
				Name: v.Name,
				Expr: right,
			}, nil
		}
		return nil, p.newError("Can only assign to variables.")
	}

	return expr, nil
}

func (p *parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR, XOR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &ExprLogical{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ExprLogical{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(PLUS, MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(ASTERISK, SLASH, PERCENT) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ExprUnary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.call()
}

func (p *parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	if p.match(OPEN_PAREN) {
		openParen := p.previous()
		args := make([]Expr, 0)
		for p.peek().Type != CLOSE_PAREN {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
			if p.peek().Type == CLOSE_PAREN {
				break
			}
			if !p.match(COMMA) {
				return nil, p.newError("Expect ',' between arguments.")
			}
		}
		if !p.match(CLOSE_PAREN) {
			return nil, p.newError("Expect ')' after argument list.")
		}
		expr = &ExprCall{
			OpenParen: openParen,
			Callee:    expr,
			Args:      args,
		}
	}

	return expr, nil
}

func (p *parser) primary() (Expr, error) {
	if p.match(NUMBER, STRING, TRUE, FALSE) {
		return &ExprLiteral{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(IDENTIFIER) {
		return &ExprVariable{
			Name: p.previous(),
		}, nil
	}

	if p.match(OPEN_PAREN) {
		openingParen := p.previous()
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(CLOSE_PAREN) {
			return nil, p.newErrorAt("Parenthesis never closed.", openingParen)
		}
		return &ExprGrouping{
			Expr: expr,
		}, nil
	}

	return nil, p.newError(fmt.Sprintf("Unexpected token '%s'", p.peek().Lexeme))
}

func (p *parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.peek().Type == t {
			p.current++
			return true
		}
	}
	return false
}

func (p *parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *parser) peek() Token {
	return p.tokens[p.current]
}

func (p *parser) synchronize() {
	if p.peek().Type == EOF {
		return
	}
	p.current++
	for p.peek().Type != EOF {
		switch p.peek().Type {
		case SEMICOLON:
			p.current++
			return
		case VAR, FUNC, IF:
			return
		}
		p.current++
	}
}

type ParseError struct {
	Token   Token
	Message string
	Line    []rune
}

func (p ParseError) Error() string {
	return generateErrorText(p.Message, p.Line, p.Token.Line, p.Token.Column, p.Token.Column+len([]rune(p.Token.Lexeme)))
}

func (p *parser) newError(message string) error {
	return ParseError{
		Token:   p.peek(),
		Message: message,
		Line:    p.lines[p.peek().Line],
	}
}

func (p *parser) newErrorAt(message string, token Token) error {
	return ParseError{
		Token:   token,
		Message: message,
		Line:    p.lines[token.Line],
	}
}
