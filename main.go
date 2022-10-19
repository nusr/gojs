package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/nusr/gojs/call"
	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/interpreter"
)

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
		result := interpreter.Interpret(line, env)
		fmt.Println(out, result)
	}
}

func RunFile(fileName string) (any, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	env := environment.New(nil)
	call.RegisterGlobal(env)
	result := interpreter.Interpret(string(content), env)
	return result, nil
}

func main() {
	argc := len(os.Args)
	if argc == 1 {
		RunCommand(os.Stdin, os.Stdout)
	} else if argc == 2 {
		fileName := os.Args[1]
		result, err := RunFile(fileName)
		if err != nil {
			fmt.Printf("can not open file \"%s\", error: %v", fileName, err)
		} else {
			fmt.Println(result)
		}
	} else {
		fmt.Println("Usage: node [path]")
		os.Exit(64)
	}
}
