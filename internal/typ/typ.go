package typ

import (
	"time"
)

// Restriction of types that can be easily converted to and from simple string representations.
type StringRepresentable interface {
	Primitive | time.Duration | string
}

type Primitive interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | bool
}

type PrimitiveSlice interface {
	[]int | []int8 | []int16 | []int32 | []int64 | []uint | []uint8 | []uint16 | []uint32 | []uint64 | []float32 | []float64 | []bool
}
