package repl

import (
	"fmt"
	"testing"

	"github.com/nusr/gojs/environment"
)

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
			//!reflect.DeepEqual(got, tt.want)
			got := interpret(tt.source, environment.New(nil))
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("expect= %v, actual= %v", tt.want, got)
			}
		})
	}
}
