package global

import (
	"github.com/nusr/gojs/call"
	"time"
)

type Clock int64

func NewClock(val int64) call.BaseCallable {
	return Clock(val)
}

func (globalClock Clock) Size() int {
	return 0
}

func (globalClock Clock) Call(interpreter call.InterpreterMethods, params []any) any {
	return time.Now().UnixMilli()
}

func (globalClock Clock) String() string {
	return "<native fn>"
}
