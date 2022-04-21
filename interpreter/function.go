package interpreter

type Callable interface {
	ArgumentCount() int
	Call(i *interpreter, args []any) (any, error)
}

type function struct {
	name       Token
	body       Stmt
	closure    *Environment
	parameters []string
}

func (f function) ArgumentCount() int {
	return len(f.parameters)
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

	return nil, err
}
