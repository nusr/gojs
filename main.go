package main

import (
	"fmt"
	"os"

	"github.com/nusr/gojs/repl"
)

func main() {
	argc := len(os.Args)
	if argc == 1 {
		repl.RunCommand(os.Stdin, os.Stdout)
	} else if argc == 2 {
		fileName := os.Args[1]
		result, err := repl.RunFile(fileName)
		if err != nil {
			fmt.Printf("can not open file \"%s\", error: %v", fileName, err)
		} else {
			fmt.Println(result)
		}
	} else {
		fmt.Println("Usage: lox [path]")
		os.Exit(64)
	}
}
