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

func Interpret(program []Stmt, lines [][]rune) error {
	interpreter := &interpreter{
		lines: lines,
		env:   NewEnvironment(nil),
	}

	for _, stmt := range program {
		err := stmt.Accept(interpreter)
		if err != nil {
			return err
		}
	}

	main, err := interpreter.env.Get("main")
	if err != nil {
		if err == ErrUndefined {
			return errors.New("No main function.")
		}
		return err
	}
	mainFunc, ok := main.(function)
	if !ok {
		return errors.New("No main function.")
	}

	return mainFunc.Call(interpreter)
}

func (i *interpreter) VisitVarDecl(stmt StmtVarDecl) error {
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

func (i *interpreter) VisitFuncDecl(stmt StmtFuncDecl) error {
	err := i.env.Define(stmt.Name.Lexeme, function{
		name: stmt.Name,
		body: stmt.Body,
	})
	if err != nil {
		if err == ErrAlreadyDefined {
			return i.newError(fmt.Sprintf("'%s' is already defined in this scope", stmt.Name.Lexeme), stmt.Name)
		}
		return i.newError(err.Error(), stmt.Name)
	}
	return nil
}

func (i *interpreter) VisitBlock(stmt StmtBlock) error {
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

func (i *interpreter) VisitLiteral(expr ExprLiteral) (any, error) {
	return expr.Value, nil
}

func (i *interpreter) VisitGrouping(expr ExprGrouping) (any, error) {
	return expr.Expr.Accept(i)
}

func (i *interpreter) VisitUnary(expr ExprUnary) (any, error) {
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

func (i *interpreter) VisitBinary(expr ExprBinary) (any, error) {
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

func (i *interpreter) VisitLogical(expr ExprLogical) (any, error) {
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
