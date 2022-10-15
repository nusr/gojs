package types

import (
	"github.com/nusr/gojs/statement"
)

type Interpreter interface {
	statement.StatementVisitor
	statement.ExpressionVisitor
	Interpret(list []statement.Statement) any
	GetGlobal() (globals Environment)
	Execute(statement statement.Statement) any
	Evaluate(expression statement.Expression) any
	ExecuteBlock(statement statement.BlockStatement, environment Environment) (result any)
}
