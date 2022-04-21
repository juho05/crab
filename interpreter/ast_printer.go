package interpreter

import (
	"fmt"
)

type ASTPrinter struct{}

type PrinterResult string

func (p PrinterResult) Error() string {
	return string(p)
}

func PrintAST(program Stmt) string {
	result := program.Accept(ASTPrinter{})
	return result.Error()
}

func (a ASTPrinter) VisitVarDecl(stmt StmtVarDecl) error {
	var expr any
	if stmt.Expr != nil {
		expr, _ = stmt.Expr.Accept(a)
	} else {
		expr = toString(nil)
	}
	return PrinterResult(fmt.Sprintf("[va] var %s = %v;", stmt.Name.Lexeme, expr))
}

func (a ASTPrinter) VisitFuncDecl(stmt StmtFuncDecl) error {
	body := stmt.Body.Accept(a)
	return PrinterResult(fmt.Sprintf("[fn] fun %s() %s", stmt.Name.Lexeme, body))
}

func (a ASTPrinter) VisitBlock(stmt StmtBlock) error {
	str := fmt.Sprintf("{\n")
	for _, s := range stmt.Statements {
		str = fmt.Sprintf("%s%v\n", str, s.Accept(a))
	}

	return PrinterResult(fmt.Sprintf("%s}", str))
}

func (a ASTPrinter) VisitLiteral(literal ExprLiteral) (any, error) {
	return toString(literal.Value), nil
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

func toString(value any) string {
	if _, ok := value.(string); ok {
		return fmt.Sprintf("\"%v\"", value)
	}
	if value == nil {
		return fmt.Sprintf("null")
	}
	return fmt.Sprintf("%v", value)
}