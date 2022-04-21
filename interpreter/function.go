package interpreter

type Callable interface {
	Call(i *interpreter)
}

type function struct {
	name    Token
	body    Stmt
	closure *Environment
}

func (f function) Call(i *interpreter) error {
	prevEnv := i.env

	i.env = f.closure
	err := f.body.Accept(i)
	i.env = prevEnv

	return err
}
