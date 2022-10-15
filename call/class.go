package call

import (
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/types"
)

type Instance struct {
	value map[any]any
}

func NewInstance() types.Property {
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

func (instance *Instance) Has(key any) bool {
	if _, ok := instance.value[key]; ok {
		return true
	}
	return false
}

type Class struct {
	methods []statement.Statement
	value   map[any]any
}

func NewClass(methods []statement.Statement) types.ClassType {
	return &Class{
		methods: methods,
		value:   make(map[any]any),
	}
}

func (class *Class) SetMethods(methods []statement.Statement) {
	class.methods = methods
}

func (class *Class) Call(interpreter types.InterpreterMethods, params []any) any {
	env := environment.New(interpreter.GetGlobal())
	instance := NewInstance()
	env.Define("this", instance)
	for _, item := range class.methods {
		if val, ok := item.(statement.FunctionStatement); ok {
			if val.Name.Lexeme == "constructor" {
				t := NewFunction(val.Body, val.Params, env)
				t.Call(interpreter, params)
			} else {
				instance.Set(val.Name.Lexeme, NewFunction(val.Body, val.Params, env))
			}
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

func (instance *Class) Get(key any) any {
	if val, ok := instance.value[key]; ok {
		return val
	}
	return nil
}

func (instance *Class) Set(key any, value any) {
	instance.value[key] = value
}
