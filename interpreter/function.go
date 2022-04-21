package interpreter

type Callable interface {
	ArgumentCount() int
	Call(i *interpreter, args []any) (any, error)
}

type function struct {
	name    Token
	body    Stmt
	closure *Environment
}

func (f function) ArgumentCount() int {
	return 0
}

func (f function) Call(i *interpreter, args []any) (any, error) {
	prevEnv := i.env

	i.env = f.closure
	err := f.body.Accept(i)
	i.env = prevEnv

	return nil, err
}
