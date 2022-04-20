package interpreter

import (
	"fmt"
	"strings"
)

type ParseError struct {
	Token   Token
	Message string
	Line    []rune
}

func (p ParseError) Error() string {
	line := p.Line
	if p.Token.Column+len([]rune(p.Token.Lexeme)) >= len(p.Line) {
		line = append(line, []rune(strings.Repeat(" ", p.Token.Column+len([]rune(p.Token.Lexeme))-(len(p.Line)-1)))...)
	}

	length := len(line)
	line = []rune(strings.TrimPrefix(strings.TrimPrefix(string(line), " "), "\t"))
	columnStart := p.Token.Column - (length - len(line))
	columnEnd := columnStart + len([]rune(p.Token.Lexeme))

	errorLine := string(line[:columnStart])
	errorLine = errorLine + "\x1b[4m\x1b[31m"
	errorLine = errorLine + string(line[columnStart:columnEnd])
	errorLine = errorLine + "\x1b[0m"
	errorLine = errorLine + string(line[columnEnd:])

	text := fmt.Sprintf("\x1b[2m[%d]  \x1b[0m%s", p.Token.Line, errorLine)
	text = fmt.Sprintf("%s%s\n%s\n%s", fmt.Sprintf("[%d:%d]: %s\n", p.Token.Line+1, p.Token.Column+1, p.Message), strings.Repeat("-", 30), text, strings.Repeat("-", 30))
	return text
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

type parser struct {
	tokens  []Token
	current int
	lines   [][]rune
}

func Parse(tokens []Token, lines [][]rune) (Expr, error) {
	parser := &parser{
		tokens: tokens,
		lines:  lines,
	}
	return parser.parse()
}

func (p *parser) parse() (Expr, error) {
	return p.expression()
}

func (p *parser) expression() (Expr, error) {
	return p.term()
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
	if p.match(MINUS) {
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
	if p.match(NUMBER) {
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
