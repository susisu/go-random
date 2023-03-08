package random

import "math/rand"

// Generator is an abstract random number generator that yields uint64 values.
type Generator interface {
	Uint64() uint64
}

var _ Generator = (rand.Source64)(nil)
