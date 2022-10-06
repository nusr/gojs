package scanner

import (
	"testing"

	"github.com/nusr/gojs/token"
)

func TestScanner(t *testing.T) {
	source := `
	'str'
	1
	1.0
	true
	false
	console.log(null);
	`
	scanner := New(source)
	tokens := scanner.Scan()

	expects := []struct {
		Type  token.Type
		Value string
	}{
		{
			token.String,
			"str",
		},
		{
			token.Int64,
			"1",
		},
		{
			token.Float64,
			"1.0",
		},
		{
			token.True,
			"true",
		},
		{
			token.False,
			"false",
		},
		{
			token.Identifier,
			"console",
		},
		{
			token.Dot,
			".",
		},
		{
			token.Identifier,
			"log",
		},
		{
			token.LeftParen,
			"(",
		},
		{
			token.Null,
			"null",
		},
		{
			token.RightParen,
			")",
		},
		{
			token.Semicolon,
			";",
		},
		{
			token.EOF,
			"",
		},
	}
	for i, item := range tokens {
		expect := expects[i]
		if item.Type != expect.Type {
			t.Errorf("token type expect= %v, actual= %v", expect.Type, item.Type)
		}
		if item.Lexeme != expect.Value {
			t.Errorf("token value expect= %v, actual= %v", expect.Value, item.Lexeme)
		}
	}
}
