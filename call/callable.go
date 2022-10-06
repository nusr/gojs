package call

import (
	"fmt"

	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
)

type InterpreterMethods interface {
	GetGlobal() (globals *environment.Environment)
	ExecuteBlock(statement statement.BlockStatement, environment *environment.Environment) (result any)
}

type BaseCallable interface {
	Call(interpreter InterpreterMethods, params []any) any
	String() string
}

type Callable struct {
	declaration statement.FunctionStatement
}

func NewCallable(declaration statement.FunctionStatement) BaseCallable {
	return &Callable{
		declaration: declaration,
	}
}

func (callable *Callable) Call(interpreter InterpreterMethods, params []any) any {
	env := environment.New(interpreter.GetGlobal())
	paramsLen := len(params)
	for i, item := range callable.declaration.Params {
		if i <= (paramsLen - 1) {
			env.Define(item.Lexeme, params[i])
		} else {
			env.Define(item.Lexeme, nil)
		}
	}
	return interpreter.ExecuteBlock(callable.declaration.Body, env)
}

func (callable *Callable) String() string {
	return fmt.Sprintf("<Fun %s>", callable.declaration.Name.Lexeme)
}
