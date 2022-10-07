package call

import (
	"testing"
)

func TestInstance(t *testing.T) {
	tests := []struct {
		actionType string
		index      any
		value      any
	}{
		{
			"get",
			int8(0),
			nil,
		},
		{
			"set",
			int16(100),
			1.0,
		},
		{
			"get",
			int32(40),
			nil,
		},
		{
			"set",
			int(30),
			true,
		},
		{
			"set",
			float32(2.0),
			2.0,
		},
		{
			"set",
			float64(3.6),
			nil,
		},
		{
			"get",
			int64(40),
			nil,
		},
	}
	arr := NewInstance()
	for _, item := range tests {
		if item.actionType == "set" {
			arr.Set(item.index, item.value)
			if arr.Get(item.index) != item.value {
				t.Errorf("instance.Set(%d) actual = %v, expect= %v", item.index, arr.Get(item.index), item.value)
			}
		} else {
			if arr.Get(item.index) != item.value {
				t.Errorf("instance.Get(%d) actual = %v, expect= %v", item.index, arr.Get(item.index), item.value)
			}
		}
	}
}
