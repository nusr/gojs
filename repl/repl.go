package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/global"
	"github.com/nusr/gojs/interpreter"
	"github.com/nusr/gojs/parser"
	"github.com/nusr/gojs/scanner"
)

func interpret(source string, env *environment.Environment) any {
	s := scanner.New(source)
	tokens := s.ScanTokens()

	p := parser.New(tokens)
	statements := p.Parse()

	env.Define("clock", global.Clock(time.Now().UnixMilli()))

	i := interpreter.New(env)
	return i.Interpret(statements)
}

func RunCommand(in io.Reader, out io.Writer) {
	input := bufio.NewScanner(in)
	env := environment.New(nil)
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
	return interpret(string(content), environment.New(nil)), nil
}
