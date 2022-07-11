package charm

import (
	"flag"
	"fmt"
	"time"
)

func Var[T FlagType](flags *flag.FlagSet, value *T, defaultValue T, name, usage string) FlagDefinition[T] {
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
		panic(fmt.Errorf(`unsupported type: %T`, v))
	}
	return FlagDefinition[T]{Name: name, Usage: usage, Value: value, DefaultValue: defaultValue}
}
