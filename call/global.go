package call

import (
	"fmt"

	"github.com/nusr/gojs/environment"
)

type GlobalFunction struct {
	name string
}

func NewGlobal(name string) Callable {
	return &GlobalFunction{
		name: name,
	}
}

func (g *GlobalFunction) Call(interpreter InterpreterMethods, params []any) any {
	if g.name == "console.log" {
		fmt.Println(params...)
	}
	return nil
}

func (g *GlobalFunction) String() string {
	return ""
}

func RegisterGlobal(env *environment.Environment) {
	instance := NewInstance()
	instance.Set("log", NewGlobal("console.log"))
	instance.Set("warn", NewGlobal("console.warn"))
	env.Define("console", instance)
}
