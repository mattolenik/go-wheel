package charm

import (
	"flag"
	"fmt"
	"reflect"
	"time"
)

type FlagType interface {
	float64 | int | int64 | uint | uint64 | string | bool | time.Duration
}

func FlagVar[T FlagType](c *Command, value *T, dfltValue T, name, usage string) *FlagDefinition[T] {
	v := Var(c.FlagSet, value, dfltValue, name, usage)
	c.Flags = append(c.Flags, v)
	return v
}

type Flag interface {
	GetName() string
	GetUsage() string
	GetDefault() any
	GetValue() any
	SetValue(any)
}

type FlagDefinition[T FlagType] struct {
	Name    string
	Usage   string
	Default T
	Value   *T
}

func (fd *FlagDefinition[T]) GetName() string {
	return fd.Name
}

func (fd *FlagDefinition[T]) GetUsage() string {
	return fd.Name
}

func (fd *FlagDefinition[T]) GetDefault() any {
	return fd.Default
}

func (fd *FlagDefinition[T]) GetValue() any {
	return *fd.Value
}

func (fd *FlagDefinition[T]) SetValue(value any) {
	valType := reflect.TypeOf(fd.Value).Elem()
	vType := reflect.TypeOf(value)
	if vType.AssignableTo(valType) {
		v := reflect.ValueOf(fd.Value)
		v.Elem().Set(reflect.ValueOf(value))
	}
}

func FlagVars[T FlagType](c *Command, flags ...FlagDefinition[T]) {
	for _, flag := range flags {
		FlagVar(c, flag.Value, flag.Default, flag.Name, flag.Usage)
	}
}

func Var[T FlagType](flags *flag.FlagSet, value *T, defaultValue T, name, usage string) *FlagDefinition[T] {
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
	return &FlagDefinition[T]{Name: name, Usage: usage, Value: value, Default: defaultValue}
}
