package interpreter

import "errors"

var (
	ErrAlreadyDefined = errors.New("The name is already defined in this scope.")
	ErrUndefined      = errors.New("Undefined name.")
)

type Environment struct {
	parent       *Environment
	names        map[string]any
	nestingLevel int
}

func NewEnvironment(parent *Environment) *Environment {
	nestingLevel := 0
	if parent != nil {
		nestingLevel = parent.nestingLevel + 1
	}
	return &Environment{
		parent:       parent,
		names:        make(map[string]any),
		nestingLevel: nestingLevel,
	}
}

func (e *Environment) Define(name string, value any) error {
	if name == "" {
		return nil
	}
	if e.Exists(name) {
		return ErrAlreadyDefined
	}
	e.names[name] = value
	return nil
}

func (e *Environment) Assign(name string, value any, nestingLevel int) {
	env := e
	for nestingLevel != env.nestingLevel {
		env = env.parent
	}
	env.names[name] = value
}

func (e *Environment) Get(name string, nestingLevel int) any {
	env := e
	for nestingLevel != env.nestingLevel {
		env = env.parent
	}
	return env.names[name]
}

func (e *Environment) Exists(name string) bool {
	_, ok := e.names[name]
	return ok
}
