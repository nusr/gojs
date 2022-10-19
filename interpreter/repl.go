package interpreter

import (
	"github.com/nusr/gojs/parser"
	"github.com/nusr/gojs/scanner"
	"github.com/nusr/gojs/types"
)

func Interpret(source string, env types.Environment) any {
	s := scanner.New(source)
	tokens := s.Scan()

	p := parser.New(tokens)
	statements := p.Parse()

	i := New(env)
	return i.Interpret(statements)
}
