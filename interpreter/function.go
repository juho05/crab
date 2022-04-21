package interpreter

type Callable interface {
	Call(i *interpreter)
}

type function struct {
	name Token
	body Stmt
}

func (f function) Call(i *interpreter) error {
	return f.body.Accept(i)
}
