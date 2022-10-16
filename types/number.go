package types

type NaN struct {
}

func (n NaN) String() string {
	return "NaN"
}

func IsNaN(value any) bool {
	if _, ok := value.(NaN); ok {
		return true
	}
	return false
}
