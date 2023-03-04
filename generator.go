package random

import "math/rand"

type Generator interface {
	Next() uint64
}

type source64Generator struct {
	s rand.Source64
}

func FromSource64(s rand.Source64) Generator {
	return &source64Generator{s}
}

func (g *source64Generator) Next() uint64 {
	return g.s.Uint64()
}
