package call

import "github.com/nusr/gojs/types"

type arrayImpl struct {
	value []any
}

func NewArray() types.Property {
	return &arrayImpl{
		value: []any{},
	}
}

func convertAnyToInt(index any) int64 {
	switch data := index.(type) {
	case int8:
		return int64(data)
	case int16:
		return int64(data)
	case int32:
		return int64(data)
	case int:
		return int64(data)
	case int64:
		return data
	case float32:
		return int64(data)
	case float64:
		return int64(data)
	default:
		return -1
	}
}

func (array *arrayImpl) Get(index any) any {
	i := convertAnyToInt(index)
	if i >= 0 && i <= int64(len(array.value)-1) {
		return array.value[i]
	}
	return nil
}

func (array *arrayImpl) Set(index any, value any) {
	i := convertAnyToInt(index)
	if i >= 0 && i <= int64(len(array.value)-1) {
		array.value[i] = value
	} else if i > int64(len(array.value)-1) {
		t := make([]any, i+1)
		copy(t, array.value)
		array.value = t
		array.value[i] = value
	}
}

func (array *arrayImpl) Has(index any) bool {
	i := convertAnyToInt(index)
	if i >= 0 && i <= int64(len(array.value)-1) {
		return true
	}
	return false
}
