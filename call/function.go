package call

import (
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
	"github.com/nusr/gojs/types"
)

type Function struct {
	env    types.Environment
	body   statement.BlockStatement
	params []token.Token
}

func NewFunction(body statement.BlockStatement, params []token.Token, env types.Environment) types.Callable {
	return &Function{
		body:   body,
		params: params,
		env:    env,
	}
}

func (function *Function) Call(interpreter types.InterpreterMethods, params []any) any {
	env := environment.New(function.env)
	paramsLen := len(params)
	for i, item := range function.params {
		if i <= (paramsLen - 1) {
			env.Define(item.Lexeme, params[i])
		} else {
			env.Define(item.Lexeme, nil)
		}
	}
	return interpreter.ExecuteBlock(function.body, env)
}

func (function *Function) String() string {
	return ""
}
