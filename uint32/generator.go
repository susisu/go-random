package random

// Generator is an abstract random number generator that yields uint32 values.
type Generator interface {
	Uint32() uint32
}
