package interpreter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func (e *Environment) registerNativeFunctions() {
	e.Define("print", funcPrint{})
	e.Define("println", funcPrintln{})
	e.Define("input", funcInput{})
	e.Define("millis", funcMillis{})
}

func registerNativeFunctions(m map[string]variableState) {
	m["print"] = variableStateUsed
	m["println"] = variableStateUsed
	m["input"] = variableStateUsed
	m["millis"] = variableStateUsed
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

type funcInput struct{}

func (p funcInput) ArgumentCount() int {
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

func (p funcMillis) Call(i *interpreter, args []any) (any, error) {
	return time.Now().UnixMilli(), nil
}
