package wheel

import "time"

type CommandLineSlice interface {
	[]bool | []int | []int8 | []int16 | []int32 | []int64 | []uint | []uint8 | []uint16 | []uint32 | []uint64 | []time.Duration | []string | []any
}

type CommandLinePrimitive interface {
	bool | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | time.Duration | string
}

type CommandLineType interface {
	CommandLinePrimitive | CommandLineSlice | JSON | any
}
