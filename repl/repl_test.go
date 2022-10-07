package repl

import (
	"testing"

	"github.com/nusr/gojs/environment"
)

func Test_interpret(t *testing.T) {
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
			env := environment.New(nil)
			got := interpret(tt.source, env)
			if got != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, got)
			}
		})
	}
}
