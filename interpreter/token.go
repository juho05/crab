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
	OPEN_BRACE  TokenType = "OPEN_BRACE"
	CLOSE_BRACE TokenType = "CLOSE_BRACE"

	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"

	AND TokenType = "AND"
	OR  TokenType = "OR"
	XOR TokenType = "XOR"

	NUMBER     TokenType = "NUMBER"
	STRING     TokenType = "STRING"
	IDENTIFIER TokenType = "IDENTIFIER"

	SEMICOLON TokenType = "SEMICOLON"
	COMMA     TokenType = "COMMA"

	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"
	VAR   TokenType = "VAR"
	FUNC  TokenType = "FUNC"

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
