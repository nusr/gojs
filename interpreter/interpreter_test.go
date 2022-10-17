package interpreter

import (
	"fmt"
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
	if val, ok := actual.(fmt.Stringer); ok {
		return val.String()
	}
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

func Test_interpret_binary(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"basic",
			`
			var a = 1; 
			a += 3;
			a;`,
			int64(4),
		},
		{
			"NaN",
			`
			var a = 1 - 'test'
			a
			`,
			"NaN",
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
func Test_interpret_exponentiation(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"4**2**3",
			"4**2**3",
			float64(65536),
		},
		{

			"(4**2)**3",
			"(4**2)**3",
			float64(4096),
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

func Test_interpret_bitwise(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"or",
			"1 | 2",
			int64(3),
		},
		{
			"and",
			"1 & 2",
			int64(0),
		},
		{
			"xor",
			"1 ^ 2",
			int64(3),
		},
		{
			"not",
			"~2.0",
			int64(-3),
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

func Test_interpret_logic(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"and false",
			`
			var a = 1
			var b = false
			a && b
			`,
			false,
		},
		{
			"and false",
			`
			var a = 1
			var b = true
			a && b
			`,
			true,
		},
		{
			"or true",
			`
			var a = 1
			var b = false
			a || b
			`,
			int64(1),
		},
		{
			"or false",
			`
			var a = 0
			var b = false
			a || b
			`,
			false,
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

func Test_interpret_control_flow(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   any
	}{
		{
			"if true",
			`
			var a = 1
			if (a) {
				a = true
			} else {
				a = false
			}
			a
			`,
			true,
		},
		{
			"if false",
			`
			var a = 0
			if (a) {
				a = true
			} else {
				a = false
			}
			a
			`,
			false,
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
			add(1, 2)
			`,
			int64(3),
		},
		{
			"expression",
			`var add = function(a,b){
				return a + b;
			}
			add(1,3.0)`,
			float64(4.0),
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
		{
			"change methods",
			`
		class Base {
		a = 1
		}
		var c = new Base()
		c.a = 2
		c.a
		`,
			int64(2),
		},
		{
			"static",
			`
			class Base {
				static a = 1
			}
			Base.a = 2
			Base.a
			`,
			int64(2),
		},
		{
			"expression",
			`
			var Base = class {
				a = 1.0
			}
			var b = new Base()
			b.a
			`,
			float64(1.0),
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
