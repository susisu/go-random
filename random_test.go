package random

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

type integer interface {
	int | int32 | int64 | uint | uint32 | uint64
}

type real interface {
	float32 | float64
}

func testSnapshot[T any](t *testing.T, generate func(g Generator) T) {
	var seed int64 = 0xc0ffee // fixed for snapshots
	g := FromSource64(rand.NewSource(seed).(rand.Source64))
	numSamples := 100

	seq := make([]T, 0, numSamples)
	for i := 0; i < numSamples; i++ {
		seq = append(seq, generate(g))
	}

	snaps.MatchSnapshot(t, seq)
}

func testUniformDistribution[T any](
	t *testing.T,
	numBins int,
	binIndex func(v T) int,
	testEach func(t *testing.T, seed int64, i int, v T),
	generate func(g Generator) T,
) {
	testRng := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := testRng.Int63()
	g := FromSource64(rand.NewSource(seed).(rand.Source64))
	numSamplesPerBin := 1024
	numSamples := numBins * numSamplesPerBin

	histogram := make([]int, numBins)
	for i := 0; i < numSamples; i++ {
		v := generate(g)

		testEach(t, seed, i, v)

		histogram[binIndex(v)]++
	}

	delta := 3 * math.Sqrt(float64(numSamplesPerBin))
	for i, c := range histogram {
		assert.InDeltaf(t, numSamplesPerBin, c, delta,
			"histogram(%d) = %d should be close to %d, (seed = %d)", i, c, numSamplesPerBin, seed)
	}
}

func testSmallIntegerUniformDistribution[T integer](
	t *testing.T,
	a, b T,
	generate func(g Generator) T,
) {
	numBins := int(b - a + 1)
	testUniformDistribution(
		t,
		numBins,
		func(v T) int {
			return int(v - a)
		},
		func(t *testing.T, seed int64, i int, v T) {
			assert.GreaterOrEqualf(t, v, a,
				"v(%d) = %d should be greater than or equal to %d (seed = %d)", i, v, a, seed)
			assert.LessOrEqualf(t, v, b,
				"v(%d) = %d should be less than or equal to %d (seed = %d)", i, v, b, seed)
		},
		generate,
	)
}

func testLargeIntegerUniformDistribution[T integer](
	t *testing.T,
	a, b T,
	generate func(g Generator) T,
) {
	numBins := 8
	n := float64(uint64(b) - uint64(a))
	testUniformDistribution(
		t,
		numBins,
		func(v T) int {
			nv := float64(uint64(v)-uint64(a)) / n
			i := int(math.Floor(nv * float64(numBins)))
			if i == numBins {
				i = numBins - 1
			}
			return i
		},
		func(t *testing.T, seed int64, i int, v T) {
			assert.GreaterOrEqualf(t, v, a,
				"v(%d) = %d should be greater than or equal to %d (seed = %d)", i, v, a, seed)
			assert.LessOrEqualf(t, v, b,
				"v(%d) = %d should be less than or equal to %d (seed = %d)", i, v, b, seed)
		},
		generate,
	)
}

func testRealUniformDistribution[T real](
	t *testing.T,
	a, b T,
	generate func(g Generator) T,
) {
	numBins := 8
	n := float64(b - a)
	testUniformDistribution(
		t,
		numBins,
		func(v T) int {
			nv := float64(v-a) / n
			i := int(math.Floor(nv * float64(numBins)))
			if i == numBins {
				i = numBins - 1
			}
			return i
		},
		func(t *testing.T, seed int64, i int, v T) {
			assert.GreaterOrEqualf(t, v, a,
				"v(%d) = %f should be greater than or equal to %f (seed = %d)", i, v, a, seed)
			assert.Lessf(t, v, b,
				"v(%d) = %f should be less than %f (seed = %d)", i, v, b, seed)
		},
		generate,
	)
}

func TestInt(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Int)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt, math.MaxInt, Int)
	})
}

func TestInt32(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Int32)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt32, math.MaxInt32, Int32)
	})
}

func TestInt64(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Int64)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt64, math.MaxInt64, Int64)
	})
}

func TestUint(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Uint)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxUint, Uint)
	})
}

func TestUint32(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Uint32)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxUint32, Uint32)
	})
}

func TestUint64(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Uint64)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxUint64, Uint64)
	})
}

func TestIntBetween(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g Generator) int {
			return IntBetween(g, -128, 127)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt, math.MaxInt, func(g Generator) int {
			return IntBetween(g, math.MinInt, math.MaxInt)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, -2, 5, func(g Generator) int {
			return IntBetween(g, -2, 5)
		})
	})
}

func TestInt32Between(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g Generator) int32 {
			return Int32Between(g, -128, 127)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt32, math.MaxInt32, func(g Generator) int32 {
			return Int32Between(g, math.MinInt32, math.MaxInt32)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, -2, 5, func(g Generator) int32 {
			return Int32Between(g, -2, 5)
		})
	})
}

func TestInt64Between(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g Generator) int64 {
			return Int64Between(g, -128, 127)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt64, math.MaxInt64, func(g Generator) int64 {
			return Int64Between(g, math.MinInt64, math.MaxInt64)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, -2, 5, func(g Generator) int64 {
			return Int64Between(g, -2, 5)
		})
	})
}

func TestUintBetween(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g Generator) uint {
			return UintBetween(g, 0, 256)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxInt, func(g Generator) uint {
			return UintBetween(g, 0, math.MaxInt)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, 2, 9, func(g Generator) uint {
			return UintBetween(g, 2, 9)
		})
	})
}

func TestUint32Between(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g Generator) uint32 {
			return Uint32Between(g, 0, 256)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxInt32, func(g Generator) uint32 {
			return Uint32Between(g, 0, math.MaxInt32)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, 2, 9, func(g Generator) uint32 {
			return Uint32Between(g, 2, 9)
		})
	})
}

func TestUint64Between(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g Generator) uint64 {
			return Uint64Between(g, 0, 256)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxInt64, func(g Generator) uint64 {
			return Uint64Between(g, 0, math.MaxInt64)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, 2, 9, func(g Generator) uint64 {
			return Uint64Between(g, 2, 9)
		})
	})
}

func TestFloat32(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Float32)
	})

	t.Run("distribution", func(t *testing.T) {
		testRealUniformDistribution(t, 0, 1.0, Float32)
	})
}

func TestFloat64(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Float64)
	})

	t.Run("distribution", func(t *testing.T) {
		testRealUniformDistribution(t, 0, 1.0, Float64)
	})
}

func TestBool(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, Bool)
	})

	t.Run("distribution", func(t *testing.T) {
		numBins := 2
		testUniformDistribution(
			t,
			numBins,
			func(v bool) int {
				if v {
					return 1
				} else {
					return 0
				}
			},
			func(t *testing.T, seed int64, i int, v bool) {
				// nothing to test
			},
			Bool,
		)
	})
}
