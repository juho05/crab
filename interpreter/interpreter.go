package interpreter

import (
	"errors"
	"fmt"
	"math"
)

type interpreter struct {
	lines [][]rune
	env   *Environment
}

type LoopControl struct {
	Type TokenType
}

func (l LoopControl) Error() string {
	return string(l.Type)
}

func Interpret(program []Stmt, lines [][]rune) error {
	interpreter := &interpreter{
		lines: lines,
		env:   NewEnvironment(nil),
	}
	interpreter.env.registerNativeFunctions()

	for _, stmt := range program {
		err := stmt.Accept(interpreter)
		if err != nil {
			return err
		}
	}

	if !interpreter.env.Exists("main") {
		return errors.New("No main function.")
	}
	main := interpreter.env.Get("main", 0)
	mainFunc, ok := main.(function)
	if !ok || mainFunc.ArgumentCount() != 0 {
		return errors.New("No main function.")
	}

	_, err := mainFunc.Call(interpreter, nil)
	return err
}

func (i *interpreter) VisitExpression(stmt *StmtExpression) error {
	_, err := stmt.Expr.Accept(i)
	return err
}

func (i *interpreter) VisitVarDecl(stmt *StmtVarDecl) error {
	var value any
	var err error
	if stmt.Expr != nil {
		value, err = stmt.Expr.Accept(i)
		if err != nil {
			return nil
		}
	}

	err = i.env.Define(stmt.Name.Lexeme, value)
	if err != nil {
		if err == ErrAlreadyDefined {
			return i.newError(fmt.Sprintf("'%s' is already defined in this scope", stmt.Name.Lexeme), stmt.Name)
		}
		return i.newError(err.Error(), stmt.Name)
	}

	return nil
}

func (i *interpreter) VisitFuncDecl(stmt *StmtFuncDecl) error {
	err := i.env.Define(stmt.Name.Lexeme, function{
		name:       stmt.Name,
		body:       stmt.Body,
		closure:    i.env,
		parameters: stmt.Parameters,
	})
	if err != nil {
		if err == ErrAlreadyDefined {
			return i.newError(fmt.Sprintf("'%s' is already defined in this scope", stmt.Name.Lexeme), stmt.Name)
		}
		return i.newError(err.Error(), stmt.Name)
	}
	return nil
}

func (i *interpreter) VisitIf(stmt *StmtIf) error {
	condition, err := stmt.Condition.Accept(i)
	if err != nil {
		return err
	}

	if isTruthy(condition) {
		return stmt.Body.Accept(i)
	} else if stmt.ElseBody != nil {
		return stmt.ElseBody.Accept(i)
	}
	return nil
}

func (i *interpreter) VisitWhile(stmt *StmtWhile) error {
	condition, err := stmt.Condition.Accept(i)
	if err != nil {
		return err
	}

	for isTruthy(condition) {
		err = stmt.Body.Accept(i)
		loopControl, ok := err.(LoopControl)
		if ok {
			if loopControl.Type == BREAK {
				break
			}
		} else if err != nil {
			return err
		}
		condition, err = stmt.Condition.Accept(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) VisitFor(stmt *StmtFor) error {
	err := stmt.Initializer.Accept(i)
	if err != nil {
		return err
	}

	condition, err := stmt.Condition.Accept(i)
	if err != nil {
		return err
	}

	for isTruthy(condition) {
		err = stmt.Body.Accept(i)
		loopControl, ok := err.(LoopControl)
		if ok {
			if loopControl.Type == BREAK {
				break
			}
		} else if err != nil {
			return err
		}
		_, err = stmt.Increment.Accept(i)
		if err != nil {
			return err
		}
		condition, err = stmt.Condition.Accept(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) VisitLoopControl(stmt *StmtLoopControl) error {
	return LoopControl{
		Type: stmt.Keyword.Type,
	}
}

func (i *interpreter) VisitBlock(stmt *StmtBlock) error {
	i.beginScope()
	defer i.endScope()

	for _, s := range stmt.Statements {
		err := s.Accept(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *interpreter) VisitLiteral(expr *ExprLiteral) (any, error) {
	return expr.Value, nil
}

func (i *interpreter) VisitVariable(variable *ExprVariable) (any, error) {
	return i.env.Get(variable.Name.Lexeme, variable.NestingLevel), nil
}

func (i *interpreter) VisitCall(call *ExprCall) (any, error) {
	expr, err := call.Callee.Accept(i)
	if err != nil {
		return nil, err
	}
	callable, ok := expr.(Callable)
	if !ok {
		return nil, i.newError("Can only call functions.", call.OpenParen)
	}

	if callable.ArgumentCount() != -1 && callable.ArgumentCount() != len(call.Args) {
		return nil, i.newError(fmt.Sprintf("Wrong argument count. Expect %d, got %d.", callable.ArgumentCount(), len(call.Args)), call.OpenParen)
	}

	args := make([]any, len(call.Args))
	for index, a := range call.Args {
		args[index], err = a.Accept(i)
		if err != nil {
			return nil, err
		}
	}

	return callable.Call(i, args)
}

func (i *interpreter) VisitGrouping(expr *ExprGrouping) (any, error) {
	return expr.Expr.Accept(i)
}

func (i *interpreter) VisitUnary(expr *ExprUnary) (any, error) {
	right, err := expr.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case MINUS:
		if isNumber(right) {
			return -right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Operand must be a number."), expr.Operator)
	case BANG:
		return !isTruthy(right), nil
	default:
		return nil, i.newError(fmt.Sprintf("Invalid unary operator '%s'.", expr.Operator.Lexeme), expr.Operator)
	}
}

func (i *interpreter) VisitBinary(expr *ExprBinary) (any, error) {
	left, err := expr.Left.Accept(i)
	if err != nil {
		return nil, err
	}
	right, err := expr.Right.Accept(i)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case PLUS:
		if isNumber(left, right) {
			return left.(float64) + right.(float64), nil
		} else if anyString(left, right) {
			return fmt.Sprintf("%v%v", left, right), nil
		}
		return nil, i.newError(fmt.Sprintf("Operands must be either both numbers or at least one of them a string."), expr.Operator)
	case MINUS:
		if isNumber(left, right) {
			return left.(float64) - right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)
	case ASTERISK:
		if isNumber(left, right) {
			return left.(float64) * right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)
	case SLASH:
		if isNumber(left, right) {
			return left.(float64) / right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)
	case PERCENT:
		if isNumber(left, right) {
			return math.Mod(left.(float64), right.(float64)), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)

	case EQUAL_EQUAL:
		return left == right, nil
	case BANG_EQUAL:
		return left != right, nil

	case LESS:
		if isNumber(left, right) {
			return left.(float64) < right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)
	case LESS_EQUAL:
		if isNumber(left, right) {
			return left.(float64) <= right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)
	case GREATER:
		if isNumber(left, right) {
			return left.(float64) > right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)
	case GREATER_EQUAL:
		if isNumber(left, right) {
			return left.(float64) >= right.(float64), nil
		}
		return nil, i.newError(fmt.Sprintf("Both operands must be numbers."), expr.Operator)

	default:
		return nil, i.newError(fmt.Sprintf("Invalid binary operator '%s'.", expr.Operator.Lexeme), expr.Operator)
	}
}

func (i *interpreter) VisitLogical(expr *ExprLogical) (any, error) {
	left, err := expr.Left.Accept(i)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == XOR {
		right, err := expr.Right.Accept(i)
		if err != nil {
			return nil, err
		}
		return isTruthy(left) != isTruthy(right), nil
	}

	if expr.Operator.Type == OR && isTruthy(left) {
		return true, nil
	}
	if expr.Operator.Type == AND && !isTruthy(left) {
		return false, nil
	}

	right, err := expr.Right.Accept(i)
	if err != nil {
		return nil, err
	}
	return isTruthy(right), nil
}

func (i *interpreter) VisitTernary(expr *ExprTernary) (any, error) {
	left, err := expr.Left.Accept(i)
	if err != nil {
		return nil, err
	}

	if expr.Operator1.Type != QUESTION_MARK {
		return nil, i.newError(fmt.Sprintf("Invalid ternary operator '%s'.", expr.Operator1.Lexeme), expr.Operator1)
	}

	if isTruthy(left) {
		return expr.Center.Accept(i)
	}
	return expr.Right.Accept(i)
}

func (i *interpreter) VisitAssign(expr *ExprAssign) (any, error) {
	value, err := expr.Expr.Accept(i)
	if err != nil {
		return nil, err
	}
	i.env.Assign(expr.Name.Lexeme, value, expr.NestingLevel)
	return value, nil
}

func isNumber(values ...any) bool {
	for _, v := range values {
		if _, ok := v.(float64); !ok {
			return false
		}
	}
	return true
}

func anyString(values ...any) bool {
	for _, v := range values {
		if _, ok := v.(string); ok {
			return true
		}
	}
	return false
}

func isTruthy(value any) bool {
	if v, ok := value.(bool); ok {
		return v
	}

	if v, ok := value.(float64); ok {
		return v != 0
	}

	if v, ok := value.(string); ok {
		return len(v) > 0
	}

	return false
}

func (i *interpreter) beginScope() {
	i.env = NewEnvironment(i.env)
}

func (i *interpreter) endScope() {
	i.env = i.env.parent
}

type RuntimeError struct {
	Token   Token
	Message string
	Line    []rune
}

func (r RuntimeError) Error() string {
	return generateErrorText(r.Message, r.Line, r.Token.Line, r.Token.Column, r.Token.Column+len([]byte(r.Token.Lexeme)))
}

func (i *interpreter) newError(message string, token Token) error {
	return RuntimeError{
		Token:   token,
		Message: message,
		Line:    i.lines[token.Line],
	}
}
