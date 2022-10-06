package call

import (
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

type Function struct {
	env    *environment.Environment
	body   statement.BlockStatement
	params []token.Token
}

func NewFunction(body statement.BlockStatement, params []token.Token, env *environment.Environment) Callable {
	return &Function{
		body:   body,
		params: params,
		env:    env,
	}
}

func (function *Function) Call(interpreter InterpreterMethods, params []any) any {
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
