package types

import (
	"github.com/nusr/gojs/statement"
)

type Function interface {
	Call(interpreter Interpreter, params []any) any
	String() string
}

type Property interface {
	Get(key any) any
	Set(key any, value any)
}

type Class interface {
	Property
	Function
	SetMethods(methods []statement.Statement)
}
