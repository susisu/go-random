package random_test

import (
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
	random "github.com/susisu/go-random/uint64"
)

func TestMain(t *testing.M) {
	v := t.Run()
	snaps.Clean(t)
	os.Exit(v)
}

func initTestGenerator() rand.Source64 {
	testRng := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := testRng.Int63()
	g := rand.NewSource(seed).(rand.Source64)
	return g
}

type integer interface {
	int | int32 | int64 | uint | uint32 | uint64
}

type real interface {
	float32 | float64
}

func testSnapshot[T any](t *testing.T, generate func(g random.Generator) T) {
	var seed int64 = 0xc0ffee // fixed for snapshots
	g := rand.NewSource(seed).(rand.Source64)
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
	generate func(g random.Generator) T,
) {
	testRng := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed := testRng.Int63()
	g := rand.NewSource(seed).(rand.Source64)
	numSamplesPerBin := 2000
	numSamples := numBins * numSamplesPerBin

	histogram := make([]int, numBins)
	for i := 0; i < numSamples; i++ {
		v := generate(g)

		testEach(t, seed, i, v)

		histogram[binIndex(v)]++
	}

	delta := 4 * math.Sqrt(float64(numSamplesPerBin)*(1.0-1.0/float64(numBins)))
	for i, c := range histogram {
		assert.InDeltaf(t, numSamplesPerBin, c, delta,
			"histogram(%d) = %d should be close to %d, (seed = %d)", i, c, numSamplesPerBin, seed)
	}
}

func testSmallIntegerUniformDistribution[T integer](
	t *testing.T,
	a, b T,
	generate func(g random.Generator) T,
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
	generate func(g random.Generator) T,
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
	generate func(g random.Generator) T,
) {
	numBins := 8
	n := float64(b - a)
	testUniformDistribution(
		t,
		numBins,
		func(v T) int {
			nv := float64(v-a) / n
			return int(math.Floor(nv * float64(numBins)))
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
		testSnapshot(t, random.Int)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt, math.MaxInt, random.Int)
	})
}

func TestInt32(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Int32)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt32, math.MaxInt32, random.Int32)
	})
}

func TestInt64(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Int64)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt64, math.MaxInt64, random.Int64)
	})
}

func TestUint(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Uint)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxUint, random.Uint)
	})
}

func TestUint32(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Uint32)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxUint32, random.Uint32)
	})
}

func TestUint64(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Uint64)
	})

	t.Run("distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxUint64, random.Uint64)
	})
}

func TestIntBetween(t *testing.T) {
	t.Run("panics if min > max", func(t *testing.T) {
		g := initTestGenerator()
		assert.Panics(t, func() { random.IntBetween(g, -127, -128) })
	})

	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g random.Generator) int {
			return random.IntBetween(g, -128, 127)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt, math.MaxInt, func(g random.Generator) int {
			return random.IntBetween(g, math.MinInt, math.MaxInt)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, -2, 5, func(g random.Generator) int {
			return random.IntBetween(g, -2, 5)
		})
	})
}

func TestInt32Between(t *testing.T) {
	t.Run("panics if min > max", func(t *testing.T) {
		g := initTestGenerator()
		assert.Panics(t, func() { random.Int32Between(g, -127, -128) })
	})

	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g random.Generator) int32 {
			return random.Int32Between(g, -128, 127)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt32, math.MaxInt32, func(g random.Generator) int32 {
			return random.Int32Between(g, math.MinInt32, math.MaxInt32)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, -2, 5, func(g random.Generator) int32 {
			return random.Int32Between(g, -2, 5)
		})
	})
}

func TestInt64Between(t *testing.T) {
	t.Run("panics if min > max", func(t *testing.T) {
		g := initTestGenerator()
		assert.Panics(t, func() { random.Int64Between(g, -127, -128) })
	})

	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g random.Generator) int64 {
			return random.Int64Between(g, -128, 127)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, math.MinInt64, math.MaxInt64, func(g random.Generator) int64 {
			return random.Int64Between(g, math.MinInt64, math.MaxInt64)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, -2, 5, func(g random.Generator) int64 {
			return random.Int64Between(g, -2, 5)
		})
	})
}

func TestUintBetween(t *testing.T) {
	t.Run("panics if min > max", func(t *testing.T) {
		g := initTestGenerator()
		assert.Panics(t, func() { random.UintBetween(g, 128, 127) })
	})

	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g random.Generator) uint {
			return random.UintBetween(g, 0, 256)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxInt, func(g random.Generator) uint {
			return random.UintBetween(g, 0, math.MaxInt)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, 2, 9, func(g random.Generator) uint {
			return random.UintBetween(g, 2, 9)
		})
	})
}

func TestUint32Between(t *testing.T) {
	t.Run("panics if min > max", func(t *testing.T) {
		g := initTestGenerator()
		assert.Panics(t, func() { random.Uint32Between(g, 128, 127) })
	})

	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g random.Generator) uint32 {
			return random.Uint32Between(g, 0, 256)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxInt32, func(g random.Generator) uint32 {
			return random.Uint32Between(g, 0, math.MaxInt32)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, 2, 9, func(g random.Generator) uint32 {
			return random.Uint32Between(g, 2, 9)
		})
	})
}

func TestUint64Between(t *testing.T) {
	t.Run("panics if min > max", func(t *testing.T) {
		g := initTestGenerator()
		assert.Panics(t, func() { random.Uint64Between(g, 128, 127) })
	})

	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, func(g random.Generator) uint64 {
			return random.Uint64Between(g, 0, 256)
		})
	})

	t.Run("large distribution", func(t *testing.T) {
		testLargeIntegerUniformDistribution(t, 0, math.MaxInt64, func(g random.Generator) uint64 {
			return random.Uint64Between(g, 0, math.MaxInt64)
		})
	})

	t.Run("small distribution", func(t *testing.T) {
		testSmallIntegerUniformDistribution(t, 2, 9, func(g random.Generator) uint64 {
			return random.Uint64Between(g, 2, 9)
		})
	})
}

func TestFloat32(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Float32)
	})

	t.Run("distribution", func(t *testing.T) {
		testRealUniformDistribution(t, 0, 1.0, random.Float32)
	})
}

func TestFloat64(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Float64)
	})

	t.Run("distribution", func(t *testing.T) {
		testRealUniformDistribution(t, 0, 1.0, random.Float64)
	})
}

func TestBool(t *testing.T) {
	t.Run("snapshot", func(t *testing.T) {
		testSnapshot(t, random.Bool)
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
			random.Bool,
		)
	})
}
