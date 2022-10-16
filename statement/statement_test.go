package statement

import (
	"testing"

	"github.com/nusr/gojs/token"
)

func TestBlockStatement(t *testing.T) {
	var statements []Statement
	statements = append(statements, VariableStatement{
		Name: token.Token{
			Type:   token.Identifier,
			Lexeme: "test",
			Line:   1,
		},
		Static: false,
		Initializer: LiteralExpression{
			Value: "true",
			Type:  token.True,
		},
	},
	)
	statements = append(statements, ClassStatement{
		Name: token.Token{
			Type:   token.Identifier,
			Lexeme: "Base",
			Line:   1,
		},
		Methods: []Statement{
			VariableStatement{
				Name: token.Token{
					Type:   token.Identifier,
					Lexeme: "property",
					Line:   1,
				},
				Static: false,
				Initializer: LiteralExpression{
					Value: "1",
					Type:  token.Int64,
				},
			},
			FunctionStatement{
				Name: token.Token{
					Type:   token.Identifier,
					Lexeme: "method",
					Line:   1,
				},
				Body: BlockStatement{
					Statements: []Statement{
						ReturnStatement{
							Value: LiteralExpression{
								Value: "1.0",
								Type:  token.Float64,
							},
						},
					},
				},
				Params: []token.Token{
					{
						Type:   token.Identifier,
						Lexeme: "a",
						Line:   1,
					},
				},
			},
		},
	})
	statements = append(statements, IfStatement{
		Condition: LiteralExpression{
			Value: "1.0",
			Type:  token.Float64,
		},
		ThenBranch: BlockStatement{
			Statements: []Statement{
				ReturnStatement{
					Value: LiteralExpression{
						Value: "test",
						Type:  token.String,
					},
				},
			},
		},
	})
	statement := BlockStatement{
		Statements: statements,
	}
	expect := "{var test=true;class Base{property=1;method(a){return 1.0;}}if(1.0){return test;}}"
	if statement.String() != expect {
		t.Errorf("expect:%s,actual:%s", expect, statement.String())
	}
}
