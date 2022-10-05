package main

import (
	"bufio"
	"fmt"
	environment2 "github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/global"
	interpreter2 "github.com/nusr/gojs/interpreter"
	parser2 "github.com/nusr/gojs/parser"
	scanner2 "github.com/nusr/gojs/scanner"
	"io/ioutil"
	"os"
	"time"
)

func interpret(source string, environment *environment2.Environment) {
	scanner := scanner2.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := parser2.NewParser(tokens)
	statements := parser.Parse()

	environment.Define("clock", global.Clock(time.Now().UnixMilli()))

	interpreter := interpreter2.NewInterpreter(environment)
	interpreter.Interpret(statements)

	scanner = nil
	parser = nil
	interpreter = nil
}

func reply() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	environment := environment2.NewEnvironment(nil)
	for input.Scan() {
		line := input.Text()
		if line == ".exit" {
			break
		}
		interpret(line, environment)
		fmt.Print("> ")
	}
	environment = nil
}

func runFile(fileName string) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("can not open file \"%s\", error: %v", fileName, err)
		return
	}
	environment := environment2.NewEnvironment(nil)
	interpret(string(content), environment)
	environment = nil
}

var filePaths []string

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
