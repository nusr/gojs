package token

import (
	"testing"
)

func TestConvertAnyToString(t *testing.T) {
	tests := []struct {
		name string
		text any
		want string
	}{
		{
			"int64",
			int64(1),
			"1",
		},
		{
			"float64",
			float64(1.0),
			"1.0000000000",
		},
		{
			"nil",
			nil,
			"null",
		},
		{
			"string",
			"test",
			"test",
		},
		{
			"bool true",
			true,
			"true",
		},
		{
			"bool false",
			false,
			"false",
		},
		{
			"array",
			[]int64{},
			"",
		},
		{
			"token",
			Token{
				Type:   String,
				Lexeme: "test",
				Line:   1,
			},
			"test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertAnyToString(tt.text); got != tt.want {
				t.Errorf("actual = %v, expect= %v", got, tt.want)
			}
		})
	}
}
