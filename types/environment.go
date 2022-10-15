package types

type Environment interface {
	Get(key string) any
	Define(name string, value any)
	Assign(key string, value any)
}
