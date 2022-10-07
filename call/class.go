package call

import (
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
)

type Instance struct {
	value map[any]any
}

func NewInstance() Property {
	return &Instance{
		value: make(map[any]any),
	}
}

func (instance *Instance) Get(key any) any {
	if val, ok := instance.value[key]; ok {
		return val
	}
	return nil
}

func (instance *Instance) Set(key any, value any) {
	instance.value[key] = value
}

type Class struct {
	StaticMethods Property
	methods       []statement.Statement
}

func NewClass(methods []statement.Statement) Callable {
	return &Class{
		StaticMethods: NewInstance(),
		methods:       methods,
	}
}

func (class *Class) Call(interpreter InterpreterMethods, params []any) any {
	env := environment.New(interpreter.GetGlobal())
	instance := NewInstance()
	for _, item := range class.methods {
		if val, ok := item.(statement.FunctionStatement); ok {
			instance.Set(val.Name.Lexeme, NewFunction(val.Body, val.Params, env))
		} else if val, ok := item.(statement.VariableStatement); ok {
			var init any
			if val.Initializer != nil {
				init = interpreter.Evaluate(val.Initializer)
			}
			instance.Set(val.Name.Lexeme, init)
		}
	}
	return instance
}

func (class *Class) String() string {
	return ""
}
