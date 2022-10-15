package call

import (
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/types"
)

type instanceImpl struct {
	value map[any]any
}

func NewInstance() types.Property {
	return &instanceImpl{
		value: make(map[any]any),
	}
}

func (instance *instanceImpl) Get(key any) any {
	if val, ok := instance.value[key]; ok {
		return val
	}
	return nil
}

func (instance *instanceImpl) Set(key any, value any) {
	instance.value[key] = value
}

func (instance *instanceImpl) Has(key any) bool {
	if _, ok := instance.value[key]; ok {
		return true
	}
	return false
}

type classImpl struct {
	methods []statement.Statement
	value   map[any]any
}

func NewClass(methods []statement.Statement) types.Class {
	return &classImpl{
		methods: methods,
		value:   make(map[any]any),
	}
}

func (class *classImpl) SetMethods(methods []statement.Statement) {
	class.methods = methods
}

func (class *classImpl) Call(interpreter types.Interpreter, params []any) any {
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

func (class *classImpl) String() string {
	return ""
}

func (instance *classImpl) Get(key any) any {
	if val, ok := instance.value[key]; ok {
		return val
	}
	return nil
}

func (instance *classImpl) Set(key any, value any) {
	instance.value[key] = value
}
