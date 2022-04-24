package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type CallError struct {
	Message string
}

func (t CallError) Error() string {
	return t.Message
}

func newTypeError(value any, expectedType string) CallError {
	var provided string
	switch value.(type) {
	case float64:
		provided = "Float"
	case string:
		provided = "String"
	case bool:
		provided = "Boolean"
	case list:
		provided = "List"
	default:
		provided = reflect.TypeOf(value).String()
	}

	return CallError{
		Message: fmt.Sprintf("Wrong type. Expected '%s', got '%s'.", expectedType, provided),
	}
}

func (e *Environment) registerNativeFunctions() {
	e.Define("print", funcPrint{})
	e.Define("println", funcPrintln{})
	e.Define("input", funcInput{})
	e.Define("millis", funcMillis{})
	e.Define("len", funcLen{})
	e.Define("append", funcAppend{})
	e.Define("concat", funcConcat{})
	e.Define("remove", funcRemove{})
}

func registerNativeFunctions(m map[string]variable) {
	m["print"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
	m["println"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
	m["input"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
	m["millis"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
	m["len"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
	m["append"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
	m["concat"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
	m["remove"] = variable{
		state:    variableStateUsed,
		nameType: nameTypeFunction,
	}
}

type funcPrint struct{}

func (p funcPrint) ArgumentCount() int {
	return -1
}

func (p funcPrint) ReturnValueCount() int {
	return 0
}

func (f funcPrint) Call(i *interpreter, args []any) (any, error) {
	fmt.Print(strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(args...), "\n"), "\r"))
	return nil, nil
}

type funcPrintln struct{}

func (p funcPrintln) ArgumentCount() int {
	return -1
}

func (p funcPrintln) ReturnValueCount() int {
	return 0
}

func (f funcPrintln) Call(i *interpreter, args []any) (any, error) {
	fmt.Println(args...)
	return nil, nil
}

type funcInput struct{}

func (p funcInput) ArgumentCount() int {
	return 1
}

func (p funcInput) ReturnValueCount() int {
	return 1
}

func (f funcInput) Call(i *interpreter, args []any) (any, error) {
	fmt.Print(args[0])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), nil
}

type funcMillis struct{}

func (p funcMillis) ArgumentCount() int {
	return 0
}

func (p funcMillis) ReturnValueCount() int {
	return 1
}

func (p funcMillis) Call(i *interpreter, args []any) (any, error) {
	return time.Now().UnixMilli(), nil
}

type funcLen struct{}

func (p funcLen) ArgumentCount() int {
	return 1
}

func (p funcLen) ReturnValueCount() int {
	return 1
}

func (p funcLen) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		return len(l), nil
	}
	if s, ok := args[0].(string); ok {
		return len(s), nil
	}
	return nil, newTypeError(args[0], "List|String")
}

type funcAppend struct{}

func (p funcAppend) ArgumentCount() int {
	return 2
}

func (p funcAppend) ReturnValueCount() int {
	return 1
}

func (p funcAppend) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		return append(l, args[1]), nil
	}
	return nil, newTypeError(args[0], "List")
}

type funcConcat struct{}

func (p funcConcat) ArgumentCount() int {
	return 2
}

func (p funcConcat) ReturnValueCount() int {
	return 1
}

func (p funcConcat) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		if l2, ok := args[1].(list); ok {
			return append(l, l2...), nil
		}
		return nil, newTypeError(args[1], "List")
	}
	return nil, newTypeError(args[0], "List")
}

type funcRemove struct{}

func (p funcRemove) ArgumentCount() int {
	return 2
}

func (p funcRemove) ReturnValueCount() int {
	return 1
}

func (p funcRemove) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		if index, ok := args[1].(float64); ok && index == float64(int(index)) {
			if int(index) >= len(l) || index < 0 {
				return nil, CallError{
					Message: "List index out of bounds.",
				}
			}
			return append(l[:int(index)], l[int(index)+1:]...), nil
		}
		return nil, newTypeError(args[1], "Integer")
	}
	return nil, newTypeError(args[0], "List")
}
