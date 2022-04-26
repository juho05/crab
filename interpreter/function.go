package interpreter

type Callable interface {
	ArgumentCount() int
	ReturnValueCount() int
	Throws() bool
	Call(i *interpreter, args []any) (any, error)
}

type Return struct {
	Values []any
}

func (r Return) Error() string {
	return "return"
}

type function struct {
	name             Token
	body             Stmt
	closure          *Environment
	parameters       []string
	returnValueCount int
	throws           bool
}

type multiValueReturn []any

func (f function) Throws() bool {
	return f.throws
}

func (f function) ArgumentCount() int {
	return len(f.parameters)
}

func (f function) ReturnValueCount() int {
	return f.returnValueCount
}

func (f function) Call(i *interpreter, args []any) (any, error) {
	prevEnv := i.env
	i.env = f.closure
	i.beginScope()
	for index, a := range f.parameters {
		i.env.Define(a, args[index])
	}

	err := f.body.Accept(i)

	i.env = prevEnv

	if ret, ok := err.(Return); ok {
		if len(ret.Values) == 0 {
			return nil, nil
		}
		if len(ret.Values) == 1 {
			return ret.Values[0], nil
		}
		return multiValueReturn(ret.Values), nil
	}

	return nil, err
}
