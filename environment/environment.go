package environment

import (
	"fmt"
	"github.com/nusr/gojs/token"
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

func (environment *Environment) Get(name *token.Token) any {
	if val, ok := environment.values[name.Lexeme]; ok {
		return val
	}
	if environment.parent != nil {
		return environment.parent.Get(name)
	}
	panic(any(fmt.Sprintf("%s is not defined", name.Lexeme)))
}
func (environment *Environment) Define(name string, value any) {
	environment.values[name] = value
}

func (environment *Environment) Assign(name *token.Token, value any) {
	if _, ok := environment.values[name.Lexeme]; ok {
		environment.Define(name.Lexeme, value)
		return
	}
	if environment.parent != nil {
		environment.parent.Assign(name, value)
		return
	}
	panic(any(fmt.Sprintf("%s is not defined", name.Lexeme)))
}
