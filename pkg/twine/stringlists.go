package twine

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mattolenik/go-charm/internal/typ"
)

// TODO: Add tests for this package

// TODO: Better name
// TODO: Move to another package?
func FromDelimetedList[T typ.StringRepresentable](str, delim string) ([]T, error) {
	parts := strings.Split(str, delim)
	if len(parts) == 0 || (len(parts) == 1 || parts[0] == "") {
		return []T{}, nil
	}
	result := make([]T, len(parts))
	resultValue := reflect.ValueOf(result)
	for i, part := range parts {
		part = strings.TrimSpace(part)
		v, err := PrimitiveFromString[T](part)
		if err != nil {
			return nil, fmt.Errorf("list item %q is not a valid value: %w", part, err)
		}
		resultValue.Index(i).Set(reflect.ValueOf(v))
	}
	return result, nil
}

func PrimitiveFromString[T typ.StringRepresentable](part string) (T, error) {
	var partVal any
	var result T
	var err error
	switch v := any(result).(type) {
	case string:
		partVal = part
	case bool:
		partVal, err = strconv.ParseBool(part)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a bool value: %w", part, err)
		}
	case time.Duration:
		partVal, err = time.ParseDuration(part)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a time.Duration value: %w", part, err)
		}
	case float32:
		f, err := strconv.ParseFloat(part, 32)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a float32 value: %w", part, err)
		}
		partVal = float32(f)
	case float64:
		partVal, err = strconv.ParseFloat(part, 64)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a float64 value: %w", part, err)
		}
	case uint:
		u, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a uint value: %w", part, err)
		}
		partVal = uint(u)
	case uint8:
		u, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a uint8 value: %w", part, err)
		}
		partVal = uint8(u)
	case uint16:
		u, err := strconv.ParseUint(part, 10, 16)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a uint16 value: %w", part, err)
		}
		partVal = uint16(u)
	case uint32:
		u, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a uint32 value: %w", part, err)
		}
		partVal = uint32(u)
	case uint64:
		partVal, err = strconv.ParseUint(part, 10, 64)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a uint64 value: %w", part, err)
		}
	case int:
		u, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a int value: %w", part, err)
		}
		partVal = int(u)
	case int8:
		u, err := strconv.ParseInt(part, 10, 8)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a int32 value: %w", part, err)
		}
		partVal = int8(u)
	case int16:
		u, err := strconv.ParseInt(part, 10, 16)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a int32 value: %w", part, err)
		}
		partVal = int16(u)
	case int32:
		u, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a int32 value: %w", part, err)
		}
		partVal = int32(u)
	case int64:
		u, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return result, fmt.Errorf("list item %q is not a int64 value: %w", part, err)
		}
		partVal = u
	default:
		return result, fmt.Errorf("unsupported type %T", v)
	}
	// Reflection required here because Go won't allow setting of a generic value from within a
	// type switch like the one above.
	// It returns an error like this:
	//     "cannot use int64(u) (value of type int64) as T value in assignment"
	// Until the Go type system supports this (if it does), reflection must be used here.
	reflect.ValueOf(&result).Elem().Set(reflect.ValueOf(partVal))
	return result, nil
}
