package repl

import (
	"testing"

	"github.com/nusr/gojs/call"
	"github.com/nusr/gojs/environment"
)

var logData any

type globalFunction struct {
	name string
}

func newGlobal(name string) call.Callable {
	return &globalFunction{
		name: name,
	}
}

func (g *globalFunction) Call(interpreter call.InterpreterMethods, params []any) any {
	if g.name == "console.log" {
		logData = params[0]
	}
	return nil
}

func (g *globalFunction) String() string {
	return ""
}

func RegisterGlobal(env *environment.Environment) {
	instance := call.NewInstance()
	instance.Set("log", newGlobal("console.log"))
	env.Define("console", instance)
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
			float64(1.0),
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
			env := environment.New(nil)
			got := interpret(tt.source, env)
			if got != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, got)
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
			"1",
			"var a = [1,2];a[0];",
			int64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := environment.New(nil)
			got := interpret(tt.source, env)
			if got != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, got)
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
			"1",
			"var a = {b: 1};a.b;",
			int64(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := environment.New(nil)
			got := interpret(tt.source, env)
			if got != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, got)
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
			"1",
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
			env := environment.New(nil)
			got := interpret(tt.source, env)
			if got != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, got)
			}
		})
	}
}
