package parser

import (
	"fmt"
	"testing"

	"github.com/nusr/gojs/expression"
	"github.com/nusr/gojs/scanner"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

func TestParser(t *testing.T) {

	source := `
	var a = 1;
	`
	s := scanner.New(source)
	tokens := s.ScanTokens()
	p := New(tokens)
	list := p.Parse()

	fmt.Println(list)
	expects := []statement.Statement{
		statement.VariableStatement{
			Name: token.Token{
				Type:   token.Identifier,
				Lexeme: "a",
				Line:   2,
			},
			Initializer: expression.LiteralExpression{
				Value:     "1",
				TokenType: token.Int64,
			},
		},
	}
	for i, item := range list {
		if item != expects[i] {
			t.Errorf("expect=%v,actual=%v", expects[i], item)
		}
	}
}
