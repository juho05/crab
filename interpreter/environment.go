package interpreter

import "errors"

var (
	ErrAlreadyDefined = errors.New("The name is already defined in this scope.")
	ErrUndefined      = errors.New("Undefined name.")
)

type Environment struct {
	parent *Environment
	names  map[string]any
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent: parent,
		names:  make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) error {
	if _, ok := e.names[name]; ok {
		return ErrAlreadyDefined
	}
	e.names[name] = value
	return nil
}

func (e *Environment) Assign(name string, value any) error {
	if _, ok := e.names[name]; !ok {
		if e.parent == nil {
			return ErrUndefined
		} else {
			return e.parent.Assign(name, value)
		}
	}
	e.names[name] = value
	return nil
}

func (e *Environment) Get(name string) (any, error) {
	value, ok := e.names[name]
	if !ok {
		if e.parent == nil {
			return nil, ErrUndefined
		} else {
			return e.parent.Get(name)
		}
	}
	return value, nil
}
