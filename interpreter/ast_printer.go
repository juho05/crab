package interpreter

import (
	"fmt"
)

type ASTPrinter struct{}

type PrinterResult string

func (p PrinterResult) Error() string {
	return string(p)
}

func PrintAST(program Expr) string {
	result, _ := program.Accept(ASTPrinter{})
	return result.(string)
}

func (a ASTPrinter) VisitLiteral(literal ExprLiteral) (any, error) {
	return literalToString(literal), nil
}

func (a ASTPrinter) VisitGrouping(grouping ExprGrouping) (any, error) {
	expr, _ := grouping.Expr.Accept(a)
	return fmt.Sprintf("%v", expr), nil
}

func (a ASTPrinter) VisitUnary(unary ExprUnary) (any, error) {
	right, _ := unary.Right.Accept(a)
	return fmt.Sprintf("(%s%v)", unary.Operator.Lexeme, right), nil
}

func (a ASTPrinter) VisitBinary(binary ExprBinary) (any, error) {
	left, _ := binary.Left.Accept(a)
	right, _ := binary.Right.Accept(a)
	return fmt.Sprintf("(%v %s %v)", left, binary.Operator.Lexeme, right), nil
}

func (a ASTPrinter) VisitLogical(logical ExprLogical) (any, error) {
	left, _ := logical.Left.Accept(a)
	right, _ := logical.Right.Accept(a)
	return fmt.Sprintf("(%v %s %v)", left, logical.Operator.Lexeme, right), nil
}

func literalToString(literal ExprLiteral) string {
	if _, ok := literal.Value.(string); ok {
		return fmt.Sprintf("\"%v\"", literal.Value)
	}
	return fmt.Sprintf("%v", literal.Value)
}
