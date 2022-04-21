package interpreter

import (
	"fmt"
	"strings"
)

func (e *Environment) registerNativeFunctions() {
	e.Define("print", funcPrint{})
	e.Define("println", funcPrintln{})
}

type funcPrint struct{}

func (p funcPrint) ArgumentCount() int {
	return -1
}

func (f funcPrint) Call(i *interpreter, args []any) (any, error) {
	fmt.Print(strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(args...), "\n"), "\r"))
	return nil, nil
}

type funcPrintln struct{}

func (p funcPrintln) ArgumentCount() int {
	return -1
}

func (f funcPrintln) Call(i *interpreter, args []any) (any, error) {
	fmt.Println(strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(args...), "\n"), "\r"))
	return nil, nil
}
