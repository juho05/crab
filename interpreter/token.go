package interpreter

import "fmt"

type TokenType string

const (
	PLUS     TokenType = "PLUS"
	MINUS    TokenType = "MINUS"
	ASTERISK TokenType = "ASTERISK"
	SLASH    TokenType = "SLASH"
	PERCENT  TokenType = "PERCENT"

	OPEN_PAREN  TokenType = "OPEN_PAREN"
	CLOSE_PAREN TokenType = "CLOSE_PAREN"

	NUMBER TokenType = "NUMBER"

	EOF TokenType = "EOF"
)

type Token struct {
	Line    int
	Column  int
	Type    TokenType
	Lexeme  string
	Literal any
}

func (t Token) String() string {
	return fmt.Sprintf("([%d:%d] %v %v %v)", t.Line+1, t.Column+1, t.Type, t.Lexeme, t.Literal)
}
