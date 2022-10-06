package call

type Array struct {
	value []any
}

func NewArray() Property {
	list := make([]any, 100)
	return &Array{
		value: list,
	}
}

func (array *Array) Get(index any) any {
	if val, ok := index.(int64); ok {
		return array.value[val]
	}
	return nil
}

func (array *Array) Set(index any, value any) {
	if val, ok := index.(int64); ok {
		array.value[val] = value
	}
}
