package interpreter

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
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
	"toLower":        funcToLower{},
	"toUpper":        funcToUpper{},
	"contains":       funcContains{},
	"indexOf":        funcIndexOf{},
	"trim":           funcTrim{},
	"replace":        funcReplace{},
	"split":          funcSplit{},
	"join":           funcJoin{},
	"random":         funcRandom{},
	"randomInt":      funcRandomInt{},
	"min":            funcMin{},
	"max":            funcMax{},
	"floor":          funcFloor{},
	"ceil":           funcCeil{},
	"round":          funcRound{},
	"sqrt":           funcSqrt{},
}

type funcPrint struct{}

func (f funcPrint) Throws() bool {
	return false
}

func (f funcPrint) ArgumentCount() int {
	return -1
}

func (f funcPrint) ReturnValueCount() int {
	return 0
}

func (f funcPrint) Call(i *interpreter, args []any) (any, error) {
	fmt.Print(strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintln(args...), "\n"), "\r"))
	return nil, nil
}

type funcPrintln struct{}

func (f funcPrintln) Throws() bool {
	return false
}

func (f funcPrintln) ArgumentCount() int {
	return -1
}

func (f funcPrintln) ReturnValueCount() int {
	return 0
}

func (f funcPrintln) Call(i *interpreter, args []any) (any, error) {
	fmt.Println(args...)
	return nil, nil
}

type funcInput struct{}

func (f funcInput) Throws() bool {
	return false
}

func (f funcInput) ArgumentCount() int {
	return 1
}

func (f funcInput) ReturnValueCount() int {
	return 1
}

func (f funcInput) Call(i *interpreter, args []any) (any, error) {
	fmt.Print(args[0])
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text(), nil
}

type funcMillis struct{}

func (f funcMillis) Throws() bool {
	return false
}

func (f funcMillis) ArgumentCount() int {
	return 0
}

func (f funcMillis) ReturnValueCount() int {
	return 1
}

func (f funcMillis) Call(i *interpreter, args []any) (any, error) {
	return time.Now().UnixMilli(), nil
}

type funcToString struct{}

func (f funcToString) Throws() bool {
	return false
}

func (f funcToString) ArgumentCount() int {
	return 1
}

func (f funcToString) ReturnValueCount() int {
	return 1
}

func (f funcToString) Call(i *interpreter, args []any) (any, error) {
	return fmt.Sprint(args[0]), nil
}

type funcToNumber struct{}

func (f funcToNumber) Throws() bool {
	return true
}

func (f funcToNumber) ArgumentCount() int {
	return 1
}

func (f funcToNumber) ReturnValueCount() int {
	return 1
}

func (f funcToNumber) Call(i *interpreter, args []any) (any, error) {
	number, err := strconv.ParseFloat(fmt.Sprint(args[0]), 64)
	if err != nil {
		return nil, i.NewException(fmt.Sprintf("Cannot convert '%v' to a number.", args[0]), -1)
	}
	return number, nil
}

type funcToBoolean struct{}

func (f funcToBoolean) Throws() bool {
	return true
}

func (f funcToBoolean) ArgumentCount() int {
	return 1
}

func (f funcToBoolean) ReturnValueCount() int {
	return 1
}

func (f funcToBoolean) Call(i *interpreter, args []any) (any, error) {
	boolean, err := strconv.ParseBool(fmt.Sprint(args[0]))
	if err != nil {
		return nil, i.NewException(fmt.Sprintf("Cannot convert '%v' to a boolean.", args[0]), -1)
	}
	return boolean, nil
}

type funcLen struct{}

func (f funcLen) Throws() bool {
	return false
}

func (f funcLen) ArgumentCount() int {
	return 1
}

func (f funcLen) ReturnValueCount() int {
	return 1
}

func (f funcLen) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		return float64(len(l)), nil
	}
	if s, ok := args[0].(string); ok {
		return float64(len(s)), nil
	}
	return nil, newTypeError(args[0], "List|String")
}

type funcAppend struct{}

func (f funcAppend) Throws() bool {
	return false
}

func (f funcAppend) ArgumentCount() int {
	return 2
}

func (f funcAppend) ReturnValueCount() int {
	return 1
}

func (f funcAppend) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		return append(l, args[1]), nil
	}
	return nil, newTypeError(args[0], "List")
}

type funcConcat struct{}

func (f funcConcat) Throws() bool {
	return false
}

func (f funcConcat) ArgumentCount() int {
	return 2
}

func (f funcConcat) ReturnValueCount() int {
	return 1
}

func (f funcConcat) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		if l2, ok := args[1].(list); ok {
			return append(l, l2...), nil
		}
		return nil, newTypeError(args[1], "List")
	}
	return nil, newTypeError(args[0], "List")
}

type funcRemove struct{}

func (f funcRemove) Throws() bool {
	return false
}

func (f funcRemove) ArgumentCount() int {
	return 2
}

func (f funcRemove) ReturnValueCount() int {
	return 1
}

func (f funcRemove) Call(i *interpreter, args []any) (any, error) {
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

func (f funcFileExists) Throws() bool {
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

func (f funcReadFileText) Throws() bool {
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

func (f funcWriteFileText) Throws() bool {
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

func (f funcAppendFileText) Throws() bool {
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

func (f funcDeleteFile) Throws() bool {
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

func (f funcListFiles) Throws() bool {
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

type funcToLower struct{}

func (f funcToLower) Throws() bool {
	return false
}

func (f funcToLower) ArgumentCount() int {
	return 1
}

func (f funcToLower) ReturnValueCount() int {
	return 1
}

func (f funcToLower) Call(i *interpreter, args []any) (any, error) {
	str := fmt.Sprint(args[0])
	return strings.ToLower(str), nil
}

type funcToUpper struct{}

func (f funcToUpper) Throws() bool {
	return false
}

func (f funcToUpper) ArgumentCount() int {
	return 1
}

func (f funcToUpper) ReturnValueCount() int {
	return 1
}

func (f funcToUpper) Call(i *interpreter, args []any) (any, error) {
	str := fmt.Sprint(args[0])
	return strings.ToUpper(str), nil
}

type funcContains struct{}

func (f funcContains) Throws() bool {
	return false
}

func (f funcContains) ArgumentCount() int {
	return 2
}

func (f funcContains) ReturnValueCount() int {
	return 1
}

func (f funcContains) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		for _, item := range l {
			if areEqual(args[1], item) {
				return true, nil
			}
		}
		return false, nil
	}

	str := fmt.Sprint(args[0])
	substring := fmt.Sprint(args[1])
	return strings.Contains(str, substring), nil
}

type funcIndexOf struct{}

func (f funcIndexOf) Throws() bool {
	return false
}

func (f funcIndexOf) ArgumentCount() int {
	return 2
}

func (f funcIndexOf) ReturnValueCount() int {
	return 1
}

func (f funcIndexOf) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		for index, item := range l {
			if areEqual(args[1], item) {
				return index, nil
			}
		}
		return -1, nil
	}

	str := fmt.Sprint(args[0])
	substring := fmt.Sprint(args[1])
	return strings.Index(str, substring), nil
}

type funcTrim struct{}

func (f funcTrim) Throws() bool {
	return false
}

func (f funcTrim) ArgumentCount() int {
	return 1
}

func (f funcTrim) ReturnValueCount() int {
	return 1
}

func (f funcTrim) Call(i *interpreter, args []any) (any, error) {
	str := fmt.Sprint(args[0])
	return strings.TrimSpace(str), nil
}

type funcReplace struct{}

func (f funcReplace) Throws() bool {
	return false
}

func (f funcReplace) ArgumentCount() int {
	return 3
}

func (f funcReplace) ReturnValueCount() int {
	return 1
}

func (f funcReplace) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		for index, item := range l {
			if areEqual(args[1], item) {
				l[index] = args[2]
			}
		}
		return l, nil
	}

	str := fmt.Sprint(args[0])
	old := fmt.Sprint(args[1])
	new := fmt.Sprint(args[2])
	return strings.ReplaceAll(str, old, new), nil
}

type funcSplit struct{}

func (f funcSplit) Throws() bool {
	return false
}

func (f funcSplit) ArgumentCount() int {
	return 2
}

func (f funcSplit) ReturnValueCount() int {
	return 1
}

func (f funcSplit) Call(i *interpreter, args []any) (any, error) {
	if l, ok := args[0].(list); ok {
		lists := make(list, 0, 1)
		segStart := 0
		for index := 0; index < len(l); index++ {
			if areEqual(args[1], l[index]) {
				lists = append(lists, l[segStart:index])
				segStart = index + 1
			}
		}
		lists = append(lists, l[segStart:])
		return lists, nil
	}

	str := fmt.Sprint(args[0])
	sep := fmt.Sprint(args[1])

	parts := strings.Split(str, sep)
	l := make(list, len(parts))
	for index, p := range parts {
		l[index] = p
	}
	return l, nil
}

type funcJoin struct{}

func (f funcJoin) Throws() bool {
	return false
}

func (f funcJoin) ArgumentCount() int {
	return 2
}

func (f funcJoin) ReturnValueCount() int {
	return 1
}

func (f funcJoin) Call(i *interpreter, args []any) (any, error) {
	l, ok := args[0].(list)
	if !ok {
		return args[0], nil
	}
	sep := fmt.Sprint(args[1])

	elems := make([]string, len(l))
	for index, item := range l {
		elems[index] = fmt.Sprint(item)
	}

	return strings.Join(elems, sep), nil
}

type funcRandom struct{}

func (f funcRandom) Throws() bool {
	return false
}

func (f funcRandom) ArgumentCount() int {
	return 2
}

func (f funcRandom) ReturnValueCount() int {
	return 1
}

func (f funcRandom) Call(i *interpreter, args []any) (any, error) {
	num1 := 0.0
	if n1, ok := args[0].(float64); ok {
		num1 = n1
	} else {
		return nil, newTypeError(args[0], "Number")
	}
	num2 := 0.0
	if n2, ok := args[1].(float64); ok {
		num2 = n2
	} else {
		return nil, newTypeError(args[1], "Number")
	}

	if num1 > num2 {
		return nil, CallError{
			Message: fmt.Sprintf("Second argument is less than the first argument."),
		}
	}

	return rand.Float64()*(num2-num1) + num1, nil
}

type funcRandomInt struct{}

func (f funcRandomInt) Throws() bool {
	return false
}

func (f funcRandomInt) ArgumentCount() int {
	return 2
}

func (f funcRandomInt) ReturnValueCount() int {
	return 1
}

func (f funcRandomInt) Call(i *interpreter, args []any) (any, error) {
	num1 := 0.0
	if n1, ok := args[0].(float64); ok && n1 == float64(int64(n1)) {
		num1 = n1
	} else {
		return nil, newTypeError(args[0], "Integer")
	}
	num2 := 0.0
	if n2, ok := args[1].(float64); ok && n2 == float64(int64(n2)) {
		num2 = n2
	} else {
		return nil, newTypeError(args[1], "Integer")
	}

	if num1 > num2 {
		return nil, CallError{
			Message: fmt.Sprintf("Second argument is less than the first argument."),
		}
	}

	return float64(int(rand.Float64()*(num2-num1) + num1)), nil
}

type funcMin struct{}

func (f funcMin) Throws() bool {
	return false
}

func (f funcMin) ArgumentCount() int {
	return 2
}

func (f funcMin) ReturnValueCount() int {
	return 1
}

func (f funcMin) Call(i *interpreter, args []any) (any, error) {
	num1 := 0.0
	if n1, ok := args[0].(float64); ok {
		num1 = n1
	} else {
		return nil, newTypeError(args[0], "Number")
	}
	num2 := 0.0
	if n2, ok := args[1].(float64); ok {
		num2 = n2
	} else {
		return nil, newTypeError(args[1], "Number")
	}

	return math.Min(num1, num2), nil
}

type funcMax struct{}

func (f funcMax) Throws() bool {
	return false
}

func (f funcMax) ArgumentCount() int {
	return 2
}

func (f funcMax) ReturnValueCount() int {
	return 1
}

func (f funcMax) Call(i *interpreter, args []any) (any, error) {
	num1 := 0.0
	if n1, ok := args[0].(float64); ok {
		num1 = n1
	} else {
		return nil, newTypeError(args[0], "Number")
	}
	num2 := 0.0
	if n2, ok := args[1].(float64); ok {
		num2 = n2
	} else {
		return nil, newTypeError(args[1], "Number")
	}

	return math.Max(num1, num2), nil
}

type funcAbs struct{}

func (f funcAbs) Throws() bool {
	return false
}

func (f funcAbs) ArgumentCount() int {
	return 1
}

func (f funcAbs) ReturnValueCount() int {
	return 1
}

func (f funcAbs) Call(i *interpreter, args []any) (any, error) {
	num := 0.0
	if n, ok := args[0].(float64); ok {
		num = n
	} else {
		return nil, newTypeError(args[0], "Number")
	}

	return math.Abs(num), nil
}

type funcFloor struct{}

func (f funcFloor) Throws() bool {
	return false
}

func (f funcFloor) ArgumentCount() int {
	return 1
}

func (f funcFloor) ReturnValueCount() int {
	return 1
}

func (f funcFloor) Call(i *interpreter, args []any) (any, error) {
	num := 0.0
	if n, ok := args[0].(float64); ok {
		num = n
	} else {
		return nil, newTypeError(args[0], "Number")
	}

	return math.Floor(num), nil
}

type funcCeil struct{}

func (f funcCeil) Throws() bool {
	return false
}

func (f funcCeil) ArgumentCount() int {
	return 1
}

func (f funcCeil) ReturnValueCount() int {
	return 1
}

func (f funcCeil) Call(i *interpreter, args []any) (any, error) {
	num := 0.0
	if n, ok := args[0].(float64); ok {
		num = n
	} else {
		return nil, newTypeError(args[0], "Number")
	}

	return math.Ceil(num), nil
}

type funcRound struct{}

func (f funcRound) Throws() bool {
	return false
}

func (f funcRound) ArgumentCount() int {
	return 1
}

func (f funcRound) ReturnValueCount() int {
	return 1
}

func (f funcRound) Call(i *interpreter, args []any) (any, error) {
	num := 0.0
	if n, ok := args[0].(float64); ok {
		num = n
	} else {
		return nil, newTypeError(args[0], "Number")
	}

	return math.Round(num), nil
}

type funcSqrt struct{}

func (f funcSqrt) Throws() bool {
	return false
}

func (f funcSqrt) ArgumentCount() int {
	return 1
}

func (f funcSqrt) ReturnValueCount() int {
	return 1
}

func (f funcSqrt) Call(i *interpreter, args []any) (any, error) {
	num := 0.0
	if n, ok := args[0].(float64); ok {
		num = n
	} else {
		return nil, newTypeError(args[0], "Number")
	}

	return math.Sqrt(num), nil
}
