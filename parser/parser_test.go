package parser

import (
	"fmt"
	"testing"

	"github.com/nusr/gojs/scanner"
	"github.com/nusr/gojs/statement"
	"github.com/nusr/gojs/token"
)

func TestParser(t *testing.T) {

	source := `
	var a = 1;
	console.log(1);
	`
	s := scanner.New(source)
	tokens := s.Scan()
	p := New(tokens)
	list := p.Parse()

	expects := []statement.Statement{
		statement.VariableStatement{
			Name: token.Token{
				Type:   token.Identifier,
				Lexeme: "a",
				Line:   2,
			},
			Initializer: statement.LiteralExpression{
				Value: "1",
				Type:  token.Int64,
			},
		},
		statement.ExpressionStatement{
			Expression: statement.CallExpression{
				Callee: statement.GetExpression{
					Object: statement.VariableExpression{
						Name: token.Token{
							Type:   token.Identifier,
							Lexeme: "console",
							Line:   3,
						},
					},
					Property: statement.VariableExpression{
						Name: token.Token{
							Type:   token.Identifier,
							Lexeme: "log",
							Line:   3,
						},
					},
				},
				Arguments: []statement.Expression{
					statement.LiteralExpression{
						Value: "1",
						Type:  token.Int64,
					},
				},
			},
		},
	}
	fmt.Println(len(expects), len(list))
	for i, item := range list {
		if item.String() != expects[i].String() {
			t.Errorf("expect=%v,actual=%v", expects[i], item)
		}
	}
}
