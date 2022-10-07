package environment

import (
	"testing"
)

func TestEnvironment(t *testing.T) {
	tests := []struct {
		actionType string
		key        string
		value      any
	}{
		{
			"get",
			"test",
			nil,
		},
		{
			"assign",
			"true",
			true,
		},
		{
			"define",
			"a",
			1.0,
		},
		{
			"assign",
			"a",
			false,
		},
	}
	env := New(New(nil))
	for _, item := range tests {
		if item.actionType == "define" {
			env.Define(item.key, item.value)
			if env.Get(item.key) != item.value {
				t.Errorf("env.Define(%s) actual = %v, expect= %v", item.key, env.Get(item.key), item.value)
			}
		} else if item.actionType == "assign" {
			env.Assign(item.key, item.value)
			if env.Get(item.key) != item.value {
				t.Errorf("env.Assign(%s) actual = %v, expect= %v", item.key, env.Get(item.key), item.value)
			}
		} else {
			if env.Get(item.key) != item.value {
				t.Errorf("env.Get(%s) actual = %v, expect= %v", item.key, env.Get(item.key), item.value)
			}
		}
	}
}
