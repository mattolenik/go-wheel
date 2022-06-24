package charm

import (
	"flag"
	"fmt"
	"time"
)

type Flag[T ArgType] struct {
	Name         string
	Usage        string
	Value        *T
	DefaultValue T
}

type ArgType interface {
	float64 | int | int64 | uint | uint64 | string | bool | time.Duration
}

func Var[T ArgType](flags *flag.FlagSet, value *T, defaultValue T, name, usage string) Flag[T] {
	switch v := any(value).(type) {
	case *int:
		flags.IntVar(v, name, any(defaultValue).(int), usage)
	case *string:
		flags.StringVar(v, name, any(defaultValue).(string), usage)
	case *bool:
		flags.BoolVar(v, name, any(defaultValue).(bool), usage)
	case *time.Duration:
		flags.DurationVar(v, name, any(defaultValue).(time.Duration), usage)
	case *float64:
		flags.Float64Var(v, name, any(defaultValue).(float64), usage)
	case *uint:
		flags.UintVar(v, name, any(defaultValue).(uint), usage)
	case *int64:
		flags.Int64Var(v, name, any(defaultValue).(int64), usage)
	case *uint64:
		flags.Uint64Var(v, name, any(defaultValue).(uint64), usage)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
	return Flag[T]{Name: name, Usage: usage, Value: value, DefaultValue: defaultValue}
}
