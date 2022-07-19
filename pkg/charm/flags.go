package charm

import (
	"fmt"
	"reflect"
	"time"

	"github.com/mattolenik/go-charm/internal/fn"
	"github.com/mattolenik/go-charm/internal/typ"
)

type FlagType interface {
	string | time.Duration | typ.Primitive | typ.PrimitiveSlice | []time.Duration | any
}

type FlagDefinition[T FlagType] struct {
	Name    string
	Usage   string
	Default T
	Value   *T
}

func (fd *FlagDefinition[T]) String() string {
	// TOOD: implement properly
	if s, ok := any(fd.Value).(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(fd.Value)
}

func FlagVars[T FlagType](c *Command, flags ...FlagDefinition[T]) {
	for _, flag := range flags {
		FlagVar(c, flag.Value, flag.Default, flag.Name, flag.Usage)
	}
}

func FlagVar[T FlagType](c *Command, value *T, defaultValue T, name, usage string) *FlagDefinition[T] {
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

	return &FlagDefinition[T]{Name: name, Usage: usage, Value: value, Default: defaultValue}
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
