package charm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/mattolenik/go-charm/internal/fn"
	"github.com/mattolenik/go-charm/internal/typ"
)

// TODO: Consider renaming types throughout this file

// FlagPrimitives are types that are natively supported by Go's flag package.
type FlagPrimitive interface {
	int | int64 | uint | uint64 | float64 | bool | string | time.Duration
}

type FlagType interface {
	FlagPrimitive | typ.PrimitiveSlice | []time.Duration | any
}

// TODO: Required flags, exclusive or paired required flags, etc
type Flag[T FlagType] struct {
	Name     string
	Usage    string
	Required bool
	Default  T
	Value    *T
}

func (f *Flag[T]) String() string {
	return fmt.Sprintf("Required: %t, Name: %s, Usage: %s, Value: %v, Default: %v", f.Required, f.Name, f.Usage, *f.Value, f.Default)
}

func (f *Flag[T]) Untyped() *Flag[any] {
	v := any(f.Value)
	return &Flag[any]{
		Name:     f.Name,
		Usage:    f.Usage,
		Required: f.Required,
		Default:  any(f.Default),
		Value:    &v,
	}
}

func FlagVars[T FlagType](c *Command, flags ...Flag[T]) {
	for _, flag := range flags {
		FlagVar(c, flag.Value, flag.Default, flag.Required, flag.Name, flag.Usage)
	}
}

// TODO: Rename after moving to flags package (maybe named options instead?)
func FlagF[T FlagType](c *Command, defaultValue T, required bool, name, usage string) *T {
	var t T
	f := FlagVar(c, &t, defaultValue, required, name, usage)
	return f.Value
}

func FlagVar[T FlagType](c *Command, value *T, defaultValue T, required bool, name, usage string) *Flag[T] {
	flags := c.flagSet
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
		t := reflect.TypeOf(v).Elem()
		cv, ok := c.FindTypeConverter(t)
		if !ok {
			panic(fmt.Errorf("no type converter for %q", t))
		}
		flagSetter, ok := cv.(FlagSetterFunc[T])
		if !ok {
			panic(fmt.Errorf("expected converter for %q to have type %q but found %q instead", t, fn.TypeOf[FlagSetterFunc[T]](), reflect.TypeOf(flagSetter)))
		}
		fv := &flagValueImpl{
			func(s string) error {
				return flagSetter(s, value)
			},
			func() string {
				return fmt.Sprintf("%v", value)
			},
		}
		flags.Var(fv, name, usage)
	}

	f := &Flag[T]{Name: name, Usage: usage, Value: value, Default: defaultValue}
	c.Flags = append(c.Flags, f.Untyped())
	return f
}

type flagValueImpl struct {
	set    func(string) error
	string func() string
}

func (f *flagValueImpl) Set(s string) error {
	return f.set(s)
}

func (f *flagValueImpl) String() string {
	return f.string()
}
