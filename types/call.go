package types

import (
	"github.com/nusr/gojs/statement"
)

type InterpreterMethods interface {
	statement.StatementVisitor
	statement.ExpressionVisitor
	Interpret(list []statement.Statement) any
	GetGlobal() (globals Environment)
	Execute(statement statement.Statement) any
	Evaluate(expression statement.Expression) any
	ExecuteBlock(statement statement.BlockStatement, environment Environment) (result any)
}

type Callable interface {
	Call(interpreter InterpreterMethods, params []any) any
	String() string
}

type Property interface {
	Get(key any) any
	Set(key any, value any)
}

type ClassType interface {
	Property
	Callable
	SetMethods(methods []statement.Statement)
}
