package interpreter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
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

var nativeFunctions = map[string]Callable{
	"print":          funcPrint{},
	"println":        funcPrintln{},
	"input":          funcInput{},
	"millis":         funcMillis{},
	"toString":       funcToString{},
	"toNumber":       funcToNumber{},
	"toBoolean":      funcToBoolean{},
	"len":            funcLen{},
	"append":         funcAppend{},
	"concat":         funcConcat{},
	"remove":         funcRemove{},
	"fileExists":     funcFileExists{},
	"readFileText":   funcReadFileText{},
	"writeFileText":  funcWriteFileText{},
	"appendFileText": funcAppendFileText{},
	"deleteFile":     funcDeleteFile{},
	"listFiles":      funcListFiles{},
}

type funcPrint struct{}

func (p funcPrint) Throws() bool {
	return false
}

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

func (p funcPrintln) Throws() bool {
	return false
}

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

func (p funcInput) Throws() bool {
	return false
}

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

func (p funcMillis) Throws() bool {
	return false
}

func (p funcMillis) ArgumentCount() int {
	return 0
}

func (p funcMillis) ReturnValueCount() int {
	return 1
}

func (p funcMillis) Call(i *interpreter, args []any) (any, error) {
	return time.Now().UnixMilli(), nil
}

type funcToString struct{}

func (p funcToString) Throws() bool {
	return false
}

func (p funcToString) ArgumentCount() int {
	return 1
}

func (p funcToString) ReturnValueCount() int {
	return 1
}

func (p funcToString) Call(i *interpreter, args []any) (any, error) {
	return fmt.Sprint(args[0]), nil
}

type funcToNumber struct{}

func (p funcToNumber) Throws() bool {
	return true
}

func (p funcToNumber) ArgumentCount() int {
	return 1
}

func (p funcToNumber) ReturnValueCount() int {
	return 1
}

func (p funcToNumber) Call(i *interpreter, args []any) (any, error) {
	number, err := strconv.ParseFloat(fmt.Sprint(args[0]), 64)
	if err != nil {
		return nil, i.NewException(fmt.Sprintf("Cannot convert '%v' to a number.", args[0]), -1)
	}
	return number, nil
}

type funcToBoolean struct{}

func (p funcToBoolean) Throws() bool {
	return true
}

func (p funcToBoolean) ArgumentCount() int {
	return 1
}

func (p funcToBoolean) ReturnValueCount() int {
	return 1
}

func (p funcToBoolean) Call(i *interpreter, args []any) (any, error) {
	boolean, err := strconv.ParseBool(fmt.Sprint(args[0]))
	if err != nil {
		return nil, i.NewException(fmt.Sprintf("Cannot convert '%v' to a boolean.", args[0]), -1)
	}
	return boolean, nil
}

type funcLen struct{}

func (p funcLen) Throws() bool {
	return false
}

func (p funcLen) ArgumentCount() int {
	return 1
}

func (p funcLen) ReturnValueCount() int {
	return 1
}

func (p funcLen) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		return float64(len(l)), nil
	}
	if s, ok := args[0].(string); ok {
		return float64(len(s)), nil
	}
	return nil, newTypeError(args[0], "List|String")
}

type funcAppend struct{}

func (p funcAppend) Throws() bool {
	return false
}

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

func (p funcConcat) Throws() bool {
	return false
}

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

func (p funcRemove) Throws() bool {
	return false
}

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

type funcFileExists struct{}

func (p funcFileExists) Throws() bool {
	return false
}

func (f funcFileExists) ArgumentCount() int {
	return 1
}

func (f funcFileExists) ReturnValueCount() int {
	return 1
}

func (f funcFileExists) Call(i *interpreter, args []any) (any, error) {
	filepath := fmt.Sprint(args[0])
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist), nil
}

type funcReadFileText struct{}

func (p funcReadFileText) Throws() bool {
	return true
}

func (f funcReadFileText) ArgumentCount() int {
	return 1
}

func (f funcReadFileText) ReturnValueCount() int {
	return 1
}

func (f funcReadFileText) Call(i *interpreter, args []any) (any, error) {
	filepath := fmt.Sprint(args[0])
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, i.NewException(err.Error(), -1)
	}
	return string(data), nil
}

type funcWriteFileText struct{}

func (p funcWriteFileText) Throws() bool {
	return true
}

func (f funcWriteFileText) ArgumentCount() int {
	return 2
}

func (f funcWriteFileText) ReturnValueCount() int {
	return 0
}

func (f funcWriteFileText) Call(i *interpreter, args []any) (any, error) {
	filepath := fmt.Sprint(args[0])
	err := os.MkdirAll(path.Dir(filepath), 0755)
	if err != nil {
		return nil, i.NewException(err.Error(), -1)
	}
	err = os.WriteFile(filepath, []byte(fmt.Sprint(args[1])), 0755)
	if err != nil {
		return nil, i.NewException(err.Error(), -1)
	}
	return nil, nil
}

type funcAppendFileText struct{}

func (p funcAppendFileText) Throws() bool {
	return true
}

func (f funcAppendFileText) ArgumentCount() int {
	return 2
}

func (f funcAppendFileText) ReturnValueCount() int {
	return 0
}

func (f funcAppendFileText) Call(i *interpreter, args []any) (any, error) {
	filepath := fmt.Sprint(args[0])
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return nil, i.NewException(err.Error(), -1)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprint(args[1]))
	if err != nil {
		return nil, i.NewException(err.Error(), -1)
	}
	return nil, nil
}

type funcDeleteFile struct{}

func (p funcDeleteFile) Throws() bool {
	return true
}

func (f funcDeleteFile) ArgumentCount() int {
	return 1
}

func (f funcDeleteFile) ReturnValueCount() int {
	return 0
}

func (f funcDeleteFile) Call(i *interpreter, args []any) (any, error) {
	filepath := fmt.Sprint(args[0])
	err := os.Remove(filepath)
	if err != nil {
		return nil, i.NewException(err.Error(), -1)
	}
	return nil, nil
}

type funcListFiles struct{}

func (p funcListFiles) Throws() bool {
	return true
}

func (f funcListFiles) ArgumentCount() int {
	return 1
}

func (f funcListFiles) ReturnValueCount() int {
	return 1
}

func (f funcListFiles) Call(i *interpreter, args []any) (any, error) {
	filepath := fmt.Sprint(args[0])
	entries, err := os.ReadDir(filepath)
	if err != nil {
		return nil, i.NewException(err.Error(), -1)
	}
	files := make(list, len(entries))
	for i, entry := range entries {
		files[i] = entry.Name()
	}
	return files, nil
}
