package random

import "math/rand"

type Generator interface {
	Uint64() uint64
}

var _ Generator = (rand.Source64)(nil)
