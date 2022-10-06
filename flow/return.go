package flow

import (
	"fmt"

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
	return fmt.Sprintf("return value : %s\n", token.ConvertAnyToString(r.Value))
}
