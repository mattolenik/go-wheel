package charm

import (
	"flag"
	"fmt"
	"time"
)

type ArgType interface {
	int | string | bool | time.Duration
}

func Var[T ArgType](flags *flag.FlagSet, value *T, defaultValue T, name, usage string) {
	switch v := any(value).(type) {
	case *int:
		flags.IntVar(v, name, any(defaultValue).(int), usage)
	case *string:
		flags.StringVar(v, name, any(defaultValue).(string), usage)
	case *bool:
		flags.BoolVar(v, name, any(defaultValue).(bool), usage)
	case *time.Duration:
		flags.DurationVar(v, name, any(defaultValue).(time.Duration), usage)
	default:
		panic(fmt.Errorf(`unhandled type: %T`, v))
	}
}
