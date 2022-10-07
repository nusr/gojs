package flow

import (
	"github.com/nusr/gojs/token"
)

type Return struct {
	Value any
}

func NewReturnValue(value any) Return {
	return Return{
		Value: value,
	}
}

func (r Return) String() string {
	return token.ConvertAnyToString(r.Value)
}
