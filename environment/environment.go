package environment

import (
	"github.com/nusr/gojs/types"
)

type environmentData struct {
	parent types.Environment
	values map[string]any
}

func New(parent types.Environment) types.Environment {
	values := make(map[string]any)
	return &environmentData{
		parent: parent,
		values: values,
	}
}

func (environment *environmentData) Get(key string) any {
	if val, ok := environment.values[key]; ok {
		return val
	}
	if environment.parent != nil {
		return environment.parent.Get(key)
	}
	return nil
}
func (environment *environmentData) Define(name string, value any) {
	environment.values[name] = value
}

func (environment *environmentData) Assign(key string, value any) {
	if _, ok := environment.values[key]; ok {
		environment.Define(key, value)
		return
	}
	if environment.parent != nil {
		environment.parent.Assign(key, value)
		return
	}
	environment.Define(key, value)
}
