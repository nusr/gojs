package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/nusr/gojs/call"
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/interpreter"
	"github.com/nusr/gojs/parser"
	"github.com/nusr/gojs/scanner"
)

func interpret(source string, env *environment.Environment) any {
	s := scanner.New(source)
	tokens := s.Scan()

	p := parser.New(tokens)
	statements := p.Parse()

	i := interpreter.New(env)
	return i.Interpret(statements)
}

func RunCommand(in io.Reader, out io.Writer) {
	input := bufio.NewScanner(in)
	env := environment.New(nil)
	call.RegisterGlobal(env)
	for {
		fmt.Fprintf(out, "> ")
		scanned := input.Scan()
		if !scanned {
			return
		}
		line := input.Text()
		result := interpret(line, env)
		fmt.Println(out, result)
	}
}

func RunFile(fileName string) (any, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	env := environment.New(nil)
	call.RegisterGlobal(env)
	result := interpret(string(content), env)
	return result, nil
}
