package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/global"
	"github.com/nusr/gojs/interpreter"
	"github.com/nusr/gojs/parser"
	"github.com/nusr/gojs/scanner"
)

func interpret(source string, environment *environment.Environment) any {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()

	p := parser.NewParser(tokens)
	statements := p.Parse()

	environment.Define("clock", global.Clock(time.Now().UnixMilli()))

	i := interpreter.NewInterpreter(environment)
	return i.Interpret(statements)
}

func reply() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	env := environment.NewEnvironment(nil)
	for input.Scan() {
		line := input.Text()
		if line == ".exit" {
			break
		}
		interpret(line, env)
		fmt.Print("> ")
	}
}

func runFile(fileName string) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("can not open file \"%s\", error: %v", fileName, err)
		return
	}
	interpret(string(content), environment.NewEnvironment(nil))
}

func main() {
	argc := len(os.Args)
	if argc == 1 {
		reply()
	} else if argc == 2 {
		runFile(os.Args[1])
	} else {
		fmt.Println("Usage: lox [path]")
		os.Exit(64)
	}
}
