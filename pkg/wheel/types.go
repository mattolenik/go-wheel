package wheel

import (
	"fmt"
	"time"
)

type CommandLineSlice interface {
	[]bool | []int | []int8 | []int16 | []int32 | []int64 | []uint | []uint8 | []uint16 | []uint32 | []uint64 | []time.Duration | []string | []any
}

type CommandLinePrimitive interface {
	bool | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | time.Duration | string
}

type CommandLineType interface {
	CommandLinePrimitive | CommandLineSlice | JSON | any
}

type UsageError struct {
	Msg string
}

func (e *UsageError) Error() string {
	return e.Msg
}

type InvalidOptionError struct {
	OptionName string
}

func (e *InvalidOptionError) Error() string {
	return fmt.Sprintf("Unrecognized option %q", e.OptionName)
}
