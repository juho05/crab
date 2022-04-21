package interpreter

import (
	"fmt"
)

type parser struct {
	tokens  []Token
	current int
	lines   [][]rune
}

func Parse(tokens []Token, lines [][]rune) ([]Stmt, []error) {
	parser := &parser{
		tokens: tokens,
		lines:  lines,
	}
	return parser.parse()
}

func (p *parser) parse() ([]Stmt, []error) {
	parseErrors := make([]error, 0)
	statements := make([]Stmt, 0)
	for p.peek().Type != EOF {
		stmt, err := p.declaration()
		if err != nil {
			parseErrors = append(parseErrors, err)
			p.synchronize()
			continue
		}
		statements = append(statements, stmt)
	}
	return statements, parseErrors
}

func (p *parser) declaration() (Stmt, error) {
	if p.match(VAR) {
		return p.varDecl()
	}

	return nil, p.newError(fmt.Sprintf("Unexpected token '%s'", p.peek().Lexeme))
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

	return StmtVarDecl{
		Name: name,
		Expr: expr,
	}, nil
}

func (p *parser) expression() (Expr, error) {
	return p.or()
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
		expr = ExprLogical{
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
		expr = ExprLogical{
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
		expr = ExprBinary{
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
		expr = ExprBinary{
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
		expr = ExprBinary{
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
		expr = ExprBinary{
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
		return ExprUnary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *parser) primary() (Expr, error) {
	if p.match(NUMBER, STRING, TRUE, FALSE) {
		return ExprLiteral{
			Value: p.previous().Literal,
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
		return ExprGrouping{
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
	p.current++
	for p.peek().Type != EOF {
		switch p.peek().Type {
		case SEMICOLON:
			p.current++
			return
		case VAR:
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

func (p parser) newError(message string) error {
	return ParseError{
		Token:   p.peek(),
		Message: message,
		Line:    p.lines[p.peek().Line],
	}
}

func (p parser) newErrorAt(message string, token Token) error {
	return ParseError{
		Token:   token,
		Message: message,
		Line:    p.lines[token.Line],
	}
}
