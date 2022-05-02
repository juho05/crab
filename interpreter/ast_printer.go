package interpreter

import (
	"fmt"
	"strings"
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

func (a ASTPrinter) VisitExpression(stmt *StmtExpression) error {
	expr, _ := stmt.Expr.Accept(a)
	return PrinterResult(fmt.Sprintf("[ex] %v;", expr))
}

func (a ASTPrinter) VisitVarDecl(stmt *StmtVarDecl) error {
	var expr any
	if stmt.Expr != nil {
		expr, _ = stmt.Expr.Accept(a)
	} else {
		expr = toString(nil)
	}
	names := stmt.Names[0].Lexeme
	for i := 1; i < len(stmt.Names); i++ {
		names = fmt.Sprintf("%s, %s", names, stmt.Names[i].Lexeme)
	}
	return PrinterResult(fmt.Sprintf("[va] var %s = %v;", names, expr))
}

func (a ASTPrinter) VisitFuncDecl(stmt *StmtFuncDecl) error {
	body := stmt.Body.Accept(a)

	throws := ""
	if stmt.Throws {
		throws = "throws"
	}

	return PrinterResult(fmt.Sprintf("[fn] func %s() %d %s %s", stmt.Name.Lexeme, stmt.ReturnValueCount, throws, body))
}

func (a ASTPrinter) VisitIf(stmt *StmtIf) error {
	condition, _ := stmt.Condition.Accept(a)

	body := stmt.Body.Accept(a).Error()
	if !strings.HasPrefix(body, "{") {
		body = fmt.Sprintf("{\n%v\n}", body)
	}

	elseBody := ""
	if stmt.ElseBody != nil {
		elseBody = stmt.ElseBody.Accept(a).Error()
		if !strings.HasPrefix(elseBody, "{") {
			elseBody = fmt.Sprintf("{\n%v\n}", elseBody)
		}
		elseBody = fmt.Sprintf("\nelse\n%s", elseBody)
	}

	return PrinterResult(fmt.Sprintf("[if] if (%v)\n%s%s", condition, body, elseBody))
}

func (a ASTPrinter) VisitTry(stmt *StmtTry) error {
	body := stmt.Body.Accept(a).Error()
	if !strings.HasPrefix(body, "{") {
		body = fmt.Sprintf("{\n%v\n}", body)
	}

	catchBody := stmt.CatchBody.Accept(a).Error()
	if !strings.HasPrefix(catchBody, "{") {
		catchBody = fmt.Sprintf("{\n%v\n}", catchBody)
	}
	if stmt.ExceptionName.Lexeme == "" {
		catchBody = fmt.Sprintf("\ncatch\n%s", catchBody)
	} else {
		catchBody = fmt.Sprintf("\ncatch (%s)\n%s", stmt.ExceptionName.Lexeme, catchBody)
	}

	return PrinterResult(fmt.Sprintf("[tr] try\n%s%s", body, catchBody))
}

func (a ASTPrinter) VisitWhile(stmt *StmtWhile) error {
	condition, _ := stmt.Condition.Accept(a)

	body := stmt.Body.Accept(a).Error()
	if !strings.HasPrefix(body, "{") {
		body = fmt.Sprintf("{\n%v\n}", body)
	}

	return PrinterResult(fmt.Sprintf("[wh] while (%v)\n%s", condition, body))
}

func (a ASTPrinter) VisitFor(stmt *StmtFor) error {
	initializer := stmt.Initializer.Accept(a)
	condition, _ := stmt.Condition.Accept(a)
	increment, _ := stmt.Increment.Accept(a)

	body := stmt.Body.Accept(a).Error()
	if !strings.HasPrefix(body, "{") {
		body = fmt.Sprintf("{\n%v\n}", body)
	}

	return PrinterResult(fmt.Sprintf("[fo] for (%v;%v;%v)\n%s", initializer, condition, increment, body))
}

func (a ASTPrinter) VisitLoopControl(stmt *StmtLoopControl) error {
	return PrinterResult(fmt.Sprintf("[lc] %s;", stmt.Keyword.Lexeme))
}

func (a ASTPrinter) VisitReturn(stmt *StmtReturn) error {
	text := "[re] return"

	for i, v := range stmt.Values {
		value, _ := v.Accept(a)
		text = fmt.Sprintf("%s %s", text, value)
		if i < len(stmt.Values)-1 {
			text = fmt.Sprintf("%s,", text)
		}
	}

	return PrinterResult(text + ";")
}

func (a ASTPrinter) VisitThrow(stmt *StmtThrow) error {
	return PrinterResult(fmt.Sprintf("throw %v;", stmt.Value))
}

func (a ASTPrinter) VisitBlock(stmt *StmtBlock) error {
	str := fmt.Sprintf("{\n")
	for _, s := range stmt.Statements {
		str = fmt.Sprintf("%s%v\n", str, s.Accept(a))
	}

	return PrinterResult(fmt.Sprintf("%s}", str))
}

func (a ASTPrinter) VisitLiteral(literal *ExprLiteral) (any, error) {
	return toString(literal.Value), nil
}

func (a ASTPrinter) VisitGrouping(grouping *ExprGrouping) (any, error) {
	expr, _ := grouping.Expr.Accept(a)
	return fmt.Sprintf("%v", expr), nil
}

func (a ASTPrinter) VisitList(list *ExprList) (any, error) {
	values := ""
	for _, value := range list.Values {
		v, _ := value.Accept(a)
		values = fmt.Sprintf("%s%v,", values, v)
	}
	values = strings.Trim(values, ",")
	return fmt.Sprintf("([%v])", values), nil
}

func (a ASTPrinter) VisitVariable(variable *ExprVariable) (any, error) {
	return fmt.Sprintf("(%s:%d)", variable.Name.Lexeme, variable.NestingLevel), nil
}

func (a ASTPrinter) VisitCall(call *ExprCall) (any, error) {
	callee, _ := call.Callee.Accept(a)
	args := ""
	for _, arg := range call.Args {
		argStr, _ := arg.Accept(a)
		args = fmt.Sprintf("%s%v,", args, argStr)
	}
	args = strings.Trim(args, ",")
	return fmt.Sprintf("(%v(%v))", callee, args), nil
}

func (a ASTPrinter) VisitSubscript(expr *ExprSubscript) (any, error) {
	object, _ := expr.Object.Accept(a)
	subscript, _ := expr.Subscript.Accept(a)
	return fmt.Sprintf("(%v[%v])", object, subscript), nil
}

func (a ASTPrinter) VisitUnary(unary *ExprUnary) (any, error) {
	right, _ := unary.Right.Accept(a)
	return fmt.Sprintf("(%s%v)", unary.Operator.Lexeme, right), nil
}

func (a ASTPrinter) VisitBinary(binary *ExprBinary) (any, error) {
	left, _ := binary.Left.Accept(a)
	right, _ := binary.Right.Accept(a)

	operator := binary.Operator.Lexeme

	if operator == "+=" || operator == "-=" || operator == "*=" || operator == "/=" || operator == "%=" || operator == "++" || operator == "--" {
		operator = string([]rune(operator)[0])
	}

	return fmt.Sprintf("(%v %s %v)", left, operator, right), nil
}

func (a ASTPrinter) VisitLogical(logical *ExprLogical) (any, error) {
	left, _ := logical.Left.Accept(a)
	right, _ := logical.Right.Accept(a)
	return fmt.Sprintf("(%v %s %v)", left, logical.Operator.Lexeme, right), nil
}

func (a ASTPrinter) VisitTernary(ternary *ExprTernary) (any, error) {
	left, _ := ternary.Left.Accept(a)
	center, _ := ternary.Center.Accept(a)
	right, _ := ternary.Right.Accept(a)
	return fmt.Sprintf("(%v %s %v %s %v)", left, ternary.Operator1.Lexeme, center, ternary.Operator2.Lexeme, right), nil
}

func (a ASTPrinter) VisitAssign(assign *ExprAssign) (any, error) {
	right, _ := assign.Expr.Accept(a)
	assignee, _ := assign.Assignees[0].Accept(a)
	assignees := fmt.Sprintf("%v", assignee)
	for i := 1; i < len(assign.Assignees); i++ {
		assignee, _ := assign.Assignees[i].Accept(a)
		assignees = fmt.Sprintf("%s, %v", assignees, assignee)
	}
	return fmt.Sprintf("(%s = %v)", assignees, right), nil
}

func (a ASTPrinter) VisitAnonymousFunction(expr *ExprAnonymousFunction) (any, error) {
	body := expr.Body.Accept(a)

	throws := ""
	if expr.Throws {
		throws = "throws"
	}

	return fmt.Sprintf("(func() %d %s %s)", expr.ReturnValueCount, throws, body), nil
}

func toString(value any) string {
	if s, ok := value.(string); ok {
		s = strings.ReplaceAll(s, "\\", "\\\\")
		s = strings.ReplaceAll(s, "\"", "\\\"")
		s = strings.ReplaceAll(s, "\r", "\\r")
		s = strings.ReplaceAll(s, "\n", "\\n")
		s = strings.ReplaceAll(s, "\t", "\\t")
		s = strings.ReplaceAll(s, "\x1b", "\\e")
		return fmt.Sprintf("\"%v\"", s)
	}
	if value == nil {
		return fmt.Sprintf("null")
	}
	return fmt.Sprintf("%v", value)
}
