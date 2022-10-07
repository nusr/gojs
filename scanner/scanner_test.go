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
	var 测试 = '测试';
	function add(a, b) {
		return a + b
	}
	add(1,2)
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
			token.Var,
			"var",
		},
		{
			token.Identifier,
			"测试",
		},
		{
			token.Equal,
			"=",
		},
		{
			token.String,
			"测试",
		},
		{
			token.Semicolon,
			";",
		},
		{
			token.Function,
			"function",
		},
		{
			token.Identifier,
			"add",
		},
		{
			token.LeftParen,
			"(",
		},
		{
			token.Identifier,
			"a",
		},
		{
			token.Comma,
			",",
		},
		{
			token.Identifier,
			"b",
		},
		{
			token.RightParen,
			")",
		},
		{
			token.LeftBrace,
			"{",
		},
		{
			token.Return,
			"return",
		},
		{
			token.Identifier,
			"a",
		},
		{
			token.Plus,
			"+",
		},
		{
			token.Identifier,
			"b",
		},
		{
			token.RightBrace,
			"}",
		},
		{
			token.Identifier,
			"add",
		},
		{
			token.LeftParen,
			"(",
		},
		{
			token.Int64,
			"1",
		},
		{
			token.Comma,
			",",
		},
		{
			token.Int64,
			"2",
		},
		{
			token.RightParen,
			")",
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
