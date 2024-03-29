package typ

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Restriction of types that can be easily converted to and from simple string representations.
type StringRepresentable interface {
	Primitive | time.Duration | string
}

type StringParsable interface {
	FromString(string) error
}

type Primitive interface {
	bool | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

type PrimitiveSlice interface {
	[]int | []int8 | []int16 | []int32 | []int64 | []uint | []uint8 | []uint16 | []uint32 | []uint64 | []float32 | []float64 | []bool
}

// TODO: rename?
// Parse generically parses a string into any simple, string-representable type, such as primitives, strings, time.Duration, etc.
func Parse[T StringRepresentable](str string) (T, error) {
	var val any
	var result T
	var err error
	switch v := any(result).(type) {
	case string:
		val = str
	case bool:
		val, err = strconv.ParseBool(str)
		if err != nil {
			return result, fmt.Errorf("string %q is not a bool value: %w", str, err)
		}
	case time.Duration:
		val, err = time.ParseDuration(str)
		if err != nil {
			return result, fmt.Errorf("string %q is not a time.Duration value: %w", str, err)
		}
	case float32:
		f, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return result, fmt.Errorf("string %q is not a float32 value: %w", str, err)
		}
		val = float32(f)
	case float64:
		val, err = strconv.ParseFloat(str, 64)
		if err != nil {
			return result, fmt.Errorf("string %q is not a float64 value: %w", str, err)
		}
	case uint:
		u, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return result, fmt.Errorf("string %q is not a uint value: %w", str, err)
		}
		val = uint(u)
	case uint8:
		u, err := strconv.ParseUint(str, 10, 8)
		if err != nil {
			return result, fmt.Errorf("string %q is not a uint8 value: %w", str, err)
		}
		val = uint8(u)
	case uint16:
		u, err := strconv.ParseUint(str, 10, 16)
		if err != nil {
			return result, fmt.Errorf("string %q is not a uint16 value: %w", str, err)
		}
		val = uint16(u)
	case uint32:
		u, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return result, fmt.Errorf("string %q is not a uint32 value: %w", str, err)
		}
		val = uint32(u)
	case uint64:
		val, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return result, fmt.Errorf("string %q is not a uint64 value: %w", str, err)
		}
	case int:
		u, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return result, fmt.Errorf("string %q is not a int value: %w", str, err)
		}
		val = int(u)
	case int8:
		u, err := strconv.ParseInt(str, 10, 8)
		if err != nil {
			return result, fmt.Errorf("string %q is not a int32 value: %w", str, err)
		}
		val = int8(u)
	case int16:
		u, err := strconv.ParseInt(str, 10, 16)
		if err != nil {
			return result, fmt.Errorf("string %q is not a int32 value: %w", str, err)
		}
		val = int16(u)
	case int32:
		u, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return result, fmt.Errorf("string %q is not a int32 value: %w", str, err)
		}
		val = int32(u)
	case int64:
		u, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return result, fmt.Errorf("string %q is not a int64 value: %w", str, err)
		}
		val = u
	default:
		z, ok := v.(StringParsable)
		if !ok {
			return result, fmt.Errorf("unsupported type %T", v)
		}
		err = z.FromString(str)
		if err != nil {
			return result, fmt.Errorf("string %q is not a %T value: %w", str, v, err)
		}
		return result, nil
	}
	// Reflection required here because Go won't allow setting of a generic value from within a
	// type switch like the one above.
	// It returns an error like this:
	//     "cannot use int64(u) (value of type int64) as T value in assignment"
	// Until the Go type system supports this (if it does), reflection must be used here.
	reflect.ValueOf(&result).Elem().Set(reflect.ValueOf(val))
	return result, nil
}
