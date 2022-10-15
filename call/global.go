package call

import (
	"fmt"

	"github.com/nusr/gojs/types"
)

type globalImpl struct {
	name string
}

func NewGlobal(name string) types.Function {
	return &globalImpl{
		name: name,
	}
}

func (g *globalImpl) Call(interpreter types.Interpreter, params []any) any {
	if g.name == "console.log" {
		fmt.Println(params...)
	}
	return nil
}

func (g *globalImpl) String() string {
	return ""
}

func RegisterGlobal(env types.Environment) {
	instance := NewInstance()
	instance.Set("log", NewGlobal("console.log"))
	instance.Set("warn", NewGlobal("console.warn"))
	env.Define("console", instance)
}
