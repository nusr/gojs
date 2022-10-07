package flow

import (
	"testing"
)

func TestReturnValue(t *testing.T) {
	r := NewReturnValue(true)
	if r.String() != "true" {
		t.Errorf("expect = true, actual=%v", r.String())
	}
}
