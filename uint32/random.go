package random

import (
	"math"
	"math/bits"
)

// Int returns a random int value within the range [-2^n, 2^n-1].
func Int(g Generator) int {
	if bits.UintSize <= 32 {
		return int(g.Uint32())
	} else {
		return int(Uint64(g))
	}
}

// Int32 returns a random int32 value within the range [-2^31, 2^31-1].
func Int32(g Generator) int32 {
	return int32(g.Uint32())
}

// Int64 returns a random int64 value within the range [-2^63, 2^63-1].
func Int64(g Generator) int64 {
	return int64(Uint64(g))
}

// Uint returns a random uint value within the range [0, 2^n-1].
func Uint(g Generator) uint {
	if bits.UintSize <= 32 {
		return uint(g.Uint32())
	} else {
		return uint(Uint64(g))
	}
}

// Uint32 returns a random uint32 value within the range [0, 2^32-1].
func Uint32(g Generator) uint32 {
	return g.Uint32()
}

// Uint64 returns a random uint64 value within the range [0, 2^64-1].
func Uint64(g Generator) uint64 {
	lo := uint64(g.Uint32())
	hi := uint64(g.Uint32())
	return (hi << 32) | lo
}

// uintAtMost returns a random uint value within the range [0, max].
func uintAtMost(g Generator, max uint) uint {
	if max == uint(math.MaxUint) {
		return Uint(g)
	} else if ((max + 1) & max) == 0 /* max like 0b11...1 */ {
		return Uint(g) & max
	} else {
		mask := uint(math.MaxUint) >> bits.LeadingZeros(max)
		for {
			v := Uint(g) & mask
			if v <= max {
				return v
			}
		}
	}
}

// uint32AtMost returns a random uint32 value within the range [0, max].
func uint32AtMost(g Generator, max uint32) uint32 {
	if max == uint32(math.MaxUint32) {
		return Uint32(g)
	} else if ((max + 1) & max) == 0 /* max like 0b11...1 */ {
		return Uint32(g) & max
	} else {
		mask := uint32(math.MaxUint32) >> bits.LeadingZeros32(max)
		for {
			v := Uint32(g) & mask
			if v <= max {
				return v
			}
		}
	}
}

// uint64AtMost returns a random uint64 value within the range [0, max].
func uint64AtMost(g Generator, max uint64) uint64 {
	if max == uint64(math.MaxUint64) {
		return Uint64(g)
	} else if ((max + 1) & max) == 0 /* max like 0b11...1 */ {
		return Uint64(g) & max
	} else {
		mask := uint64(math.MaxUint64) >> bits.LeadingZeros64(max)
		for {
			v := Uint64(g) & mask
			if v <= max {
				return v
			}
		}
	}
}

// IntBetween returns a random int value within the range [min, max].
// It panics if min > max is given.
func IntBetween(g Generator, min, max int) int {
	if min == max {
		return min
	} else if min < max {
		return int(uintAtMost(g, uint(max)-uint(min))) + min
	} else {
		panic("invalid argument to IntBetween: min must be less than or equal to ma")
	}
}

// Int32Between returns a random int32 value within the range [min, max].
// It panics if min > max is given.
func Int32Between(g Generator, min, max int32) int32 {
	if min == max {
		return min
	} else if min < max {
		return int32(uint32AtMost(g, uint32(max)-uint32(min))) + min
	} else {
		panic("invalid argument to Int32Between: min must be less than or equal to ma")
	}
}

// Int64Between returns a random int64 value within the range [min, max].
// It panics if min > max is given.
func Int64Between(g Generator, min, max int64) int64 {
	if min == max {
		return min
	} else if min < max {
		return int64(uint64AtMost(g, uint64(max)-uint64(min))) + min
	} else {
		panic("invalid argument to Int64Between: min must be less than or equal to ma")
	}
}

// UintBetween returns a random uint value within the range [min, max].
// It panics if min > max is given.
func UintBetween(g Generator, min, max uint) uint {
	if min == max {
		return min
	} else if min < max {
		return uintAtMost(g, max-min) + min
	} else {
		panic("invalid argument to UintBetween: min must be less than or equal to ma")
	}
}

// Uint32Between returns a random uint32 value within the range [min, max].
// It panics if min > max is given.
func Uint32Between(g Generator, min, max uint32) uint32 {
	if min == max {
		return min
	} else if min < max {
		return uint32AtMost(g, max-min) + min
	} else {
		panic("invalid argument to Uint32Between: min must be less than or equal to ma")
	}
}

// Uint64Between returns a random uint64 value within the range [min, max].
// It panics if min > max is given.
func Uint64Between(g Generator, min, max uint64) uint64 {
	if min == max {
		return min
	} else if min < max {
		return uint64AtMost(g, max-min) + min
	} else {
		panic("invalid argument to Uint64Between: min must be less than or equal to ma")
	}
}

// Float32 returns a random float32 value within the range [0, 1).
func Float32(g Generator) float32 {
	return float32(uint32AtMost(g, (1<<24)-1)) / (1 << 24)
}

// Float64 returns a random float64 value within the range [0, 1).
func Float64(g Generator) float64 {
	return float64(uint64AtMost(g, (1<<53)-1)) / (1 << 53)
}

// Bool returns a random bool value.
func Bool(g Generator) bool {
	return g.Uint32()&0x1 == 1
}
