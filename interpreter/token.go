package interpreter

import "fmt"

type TokenType string

const (
	PLUS           TokenType = "PLUS"
	PLUS_EQUAL     TokenType = "PLUS_EQUAL"
	PLUS_PLUS      TokenType = "PLUS_PLUS"
	MINUS          TokenType = "MINUS"
	MINUS_EQUAL    TokenType = "MINUS_EQUAL"
	MINUS_MINUS    TokenType = "MINUS_MINUS"
	ASTERISK       TokenType = "ASTERISK"
	ASTERISK_EQUAL TokenType = "ASTERISK_EQUAL"
	SLASH          TokenType = "SLASH"
	SLASH_EQUAL    TokenType = "SLASH_EQUAL"
	PERCENT        TokenType = "PERCENT"
	PERCENT_EQUAL  TokenType = "PERCENT_EQUAL"

	OPEN_PAREN    TokenType = "OPEN_PAREN"
	CLOSE_PAREN   TokenType = "CLOSE_PAREN"
	OPEN_BRACE    TokenType = "OPEN_BRACE"
	CLOSE_BRACE   TokenType = "CLOSE_BRACE"
	OPEN_BRACKET  TokenType = "OPEN_BRACKET"
	CLOSE_BRACKET TokenType = "CLOSE_BRACKET"

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

	SEMICOLON     TokenType = "SEMICOLON"
	COMMA         TokenType = "COMMA"
	QUESTION_MARK TokenType = "QUESTION_MARK"
	COLON         TokenType = "COLON"

	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	VAR      TokenType = "VAR"
	FUNC     TokenType = "FUNC"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	WHILE    TokenType = "WHILE"
	FOR      TokenType = "FOR"
	BREAK    TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"
	RETURN   TokenType = "RETURN"
	TRY      TokenType = "TRY"
	CATCH    TokenType = "CATCH"
	THROW    TokenType = "THROW"
	THROWS   TokenType = "THROWS"

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
