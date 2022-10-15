package call

import (
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
	"github.com/nusr/gojs/types"
)

type functionImpl struct {
	env    types.Environment
	body   statement.BlockStatement
	params []token.Token
}

func NewFunction(body statement.BlockStatement, params []token.Token, env types.Environment) types.Function {
	return &functionImpl{
		body:   body,
		params: params,
		env:    env,
	}
}

func (function *functionImpl) Call(interpreter types.Interpreter, params []any) any {
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

func (function *functionImpl) String() string {
	return ""
}
