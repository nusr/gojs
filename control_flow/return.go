package control_flow

import (
	"fmt"
	"github.com/nusr/gojs/token"
)

type ReturnValue struct {
	Value any
}

func NewReturnValue(value any) ReturnValue {
	return ReturnValue{
		Value: value,
	}
}

func (r ReturnValue) String() string {
	return fmt.Sprintf("control_flow value : %s\n", token.LiteralTypeToString(r.Value))
}
