package environment

import (
	"fmt"
)

type Environment struct {
	parent *Environment
	values map[string]any
}

func NewEnvironment(parent *Environment) *Environment {
	values := make(map[string]any)
	return &Environment{
		parent: parent,
		values: values,
	}
}

func (environment *Environment) Get(key string) any {
	if val, ok := environment.values[key]; ok {
		return val
	}
	if environment.parent != nil {
		return environment.parent.Get(key)
	}
	panic(any(fmt.Sprintf("%s is not defined", key)))
}
func (environment *Environment) Define(name string, value any) {
	environment.values[name] = value
}

func (environment *Environment) Assign(key string, value any) {
	if _, ok := environment.values[key]; ok {
		environment.Define(key, value)
		return
	}
	if environment.parent != nil {
		environment.parent.Assign(key, value)
		return
	}
	panic(any(fmt.Sprintf("%s is not defined", key)))
}
