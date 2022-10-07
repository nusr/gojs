package interpreter

import (
	"testing"

	"github.com/nusr/gojs/environment"
	"github.com/nusr/gojs/parser"
	"github.com/nusr/gojs/scanner"
)

func interpret(source string) any {
	env := environment.New(nil)
	s := scanner.New(source)
	tokens := s.Scan()
	p := parser.New(tokens)
	statements := p.Parse()
	i := New(env)
	actual := i.Interpret(statements)
	return actual
}

func Test_interpret_primary(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"null",
			"null",
			nil,
		},
		{
			"int",
			"1",
			int64(1),
		},
		{
			"float",
			"1.0",
			1.0,
		},
		{
			"bool true",
			"true",
			true,
		},
		{
			"bool false",
			"false",
			false,
		},
		{
			"string 1",
			"'str'",
			"str",
		},
		{
			"string 2",
			`"string"`,
			"string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual := interpret(tt.source)

			if actual != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, actual)
			}
		})
	}
}

func Test_interpret_array(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"basic",
			"var a = [1,2];a[0];",
			int64(1),
		},
		{
			"dynamic",
			"var a = [];a[1]=1;",
			int64(1),
		},
		{
			"dynamic",
			"var a = [,,];a[0];",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual := interpret(tt.source)

			if actual != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, actual)
			}
		})
	}
}

func Test_interpret_object(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"basic",
			"var a = {b:1};a.b;",
			int64(1),
		},
		{
			"dynamic",
			"var a = {};a.b=2;",
			int64(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual := interpret(tt.source)

			if actual != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, actual)
			}
		})
	}
}

func Test_interpret_function(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"basic",
			`
			function add(a, b) {
				return a + b;
			}
			add(1, 2);
			`,
			int64(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actual := interpret(tt.source)

			if actual != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, actual)
			}
		})
	}
}

func Test_interpret_class(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"basic",
			`
			var b = 'value'
			class Base {
				a = b;
			}
			var c = new Base();
			c.a;
			`,
			"value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := interpret(tt.source)

			if actual != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, actual)
			}
		})
	}
}
