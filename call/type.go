package call

import (
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
)

type InterpreterMethods interface {
	GetGlobal() (globals *environment.Environment)
	Execute(statement statement.Statement) any
	Evaluate(expression statement.Expression) any
	ExecuteBlock(statement statement.BlockStatement, environment *environment.Environment) (result any)
}

type Callable interface {
	Call(interpreter InterpreterMethods, params []any) any
	String() string
}

type Property interface {
	Get(key any) any
	Set(key any, value any)
}
