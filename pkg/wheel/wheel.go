package wheel

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mattolenik/go-charm/internal/fn"
	"github.com/mattolenik/go-charm/internal/refract"
)

type CommandLineType interface {
	bool | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | time.Duration | string | any
}

type Option interface {
	Name() string
	Usage() string
	Required() bool
	Type() reflect.Type
	Value() any
	Setter() func(string) error
}

type TypedOption[T CommandLineType] struct {
	name         string
	usage        string
	required     bool
	TypedValue   *T
	DefaultValue *T
	typ          reflect.Type
	setter       func(string) error
}

func (o *TypedOption[T]) Name() string {
	return o.name
}

func (o *TypedOption[T]) Usage() string {
	return o.usage
}

func (o *TypedOption[T]) Required() bool {
	return o.required
}

func (o *TypedOption[T]) Value() any {
	return o.TypedValue
}

func (o *TypedOption[T]) Type() reflect.Type {
	return o.typ
}

func (o *TypedOption[T]) Setter() func(string) error {
	return o.setter
}

type Command struct {
	Name        string
	Description string
	Usage       string
	Examples    []string
	Options     []Option
	SubCommands []*Command
	// TODO: make arg parsing strongly typed
	Args    []string
	parent  *Command
	Invoked bool
}

func NewCommand(name, usage, description string, examples []string) *Command {
	c := &Command{
		Name:        name,
		Usage:       usage,
		Description: description,
		Examples:    examples,
	}
	return c
}

func AddOption[T CommandLineType](c *Command, required bool, defaultValue *T, name, usage string) *TypedOption[T] {
	var value T
	o := &TypedOption[T]{
		name:         name,
		required:     required,
		usage:        usage,
		DefaultValue: defaultValue,
		TypedValue:   &value,
		typ:          refract.TypeOf[T](),
		setter:       convert(&value),
	}
	c.Options = append(c.Options, o)
	return o
}

func (c *Command) SubCommand(name, usage, description string, examples []string) *Command {
	sc := &Command{
		Name:        name,
		Usage:       usage,
		Description: description,
		Examples:    examples,
		parent:      c,
	}
	c.SubCommands = append(c.SubCommands, sc)
	return sc
}

func (c *Command) Parse(args []string) error {
	opts, remaining := ParseOptions(args)
	for opt, values := range opts {
		supportedOpts := fn.Filter(c.Options, func(o *Option) bool { return (*o).Name() == opt })
		if len(supportedOpts) == 0 {
			return fmt.Errorf("unsupported option %q, did you mean <TODO: insert help here>?", opt)
		}
		if len(supportedOpts) > 1 {
			// This is a panic rather than an error because duplicate options indicate a serious bug in the program.
			panic(fmt.Errorf("duplicate option found, %q was defined %d times, must be only once", opt, len(supportedOpts)))
		}
		o := *supportedOpts[0]
		if o.Type().Kind() == reflect.Slice {
			if len(values) == 0 {
				return fmt.Errorf("option %q requires a value", opt)
			}
			// TODO: convert values to the type needed by o.typ. Use that register converter pattern here
			continue
		}
		if len(values) == 0 {
			if o.Type() == refract.TypeOf[bool]() {
				o.Setter()("true")
				continue
			}
			return fmt.Errorf("option %q requires a value", opt)
		}
		if len(values) == 1 {
			v := values.Values()[0]
			if err := o.Setter()(v); err != nil {
				return fmt.Errorf("option %q: %w", opt, err)
			}
		}
		if len(values) > 1 {
			return fmt.Errorf("option %q can only be specified once but was found %d times", opt, len(values))
		}
	}
	if len(remaining) == 0 {
		return nil
	}
	firstRemaining := remaining[0]
	subcmd := fn.Filter(c.SubCommands, func(sc **Command) bool { return (*sc).Name == firstRemaining })
	if len(subcmd) == 0 {
		// TODO: make arg parsing strongly typed
		c.Args = remaining
		return nil
	}
	if len(subcmd) > 1 {
		panic(fmt.Errorf("duplicate command found, %q was defined %d times, must be only once", firstRemaining, len(subcmd)))
	}
	return (*subcmd[0]).Parse(remaining[1:])
}

// ParseOptions takes CLI arguments and returns a mapping of options to values, plus the remaining, non-option arguments.
// Example:
//   []string{"-a=1", "-b=2", "-c=3", "-b=4", "arg1", "arg2", "arg3"}
//     becomes:
//   map[string]Set[string]{"a":{"1"}, "b":{"2", "4"}, "c":{"3"}}, []string{"arg1", "arg2", "arg3"}
func ParseOptions(args []string) (fn.MultiMap[string, string], []string) {
	if len(args) == 0 {
		return fn.MultiMap[string, string]{}, args
	}
	flags := fn.MultiMap[string, string]{}
	var i int
	var arg string
	for i, arg = range args {
		if strings.HasPrefix(arg, "--") {
			arg = arg[2:]
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
		} else {
			return flags, args[i:]
		}
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 1 {
			flag := parts[0]
			value, ok := Index(args, i+1)
			if !ok {
				// Continue if this is the end of the list
				continue
			}
			if strings.HasPrefix(value, "-") {
				// Next arg is flag
				continue
			}
			flags.Put(flag, value)
		} else if len(parts) == 2 {
			flags.Put(parts[0], parts[1])
		} else {
			// This shouldn't be possible since strings.SplitN(2) should never return a slice of length > 2
			panic(fmt.Errorf("unexpected only 2 strings to be returned by strings.SplitN but instead got %d", len(parts)))
		}
	}
	return flags, args[i:]
}

func Index[T any](slice []T, i int) (v T, ok bool) {
	if i >= len(slice) {
		return
	}
	return slice[i], true
}

var converters = map[reflect.Type]func(string) error{}

func convert[T CommandLineType](v *T) func(string) error {
	switch any(v).(type) {
	case *bool:
		return func(s string) error {
			b, err := strconv.ParseBool(s)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetBool(b)
			return nil
		}
	case *int:
		return func(s string) error {
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetInt(int64(i))
			return nil
		}
	case *int8:
		return func(s string) error {
			i, err := strconv.ParseInt(s, 10, 8)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetInt(int64(i))
			return nil
		}
	case *int16:
		return func(s string) error {
			i, err := strconv.ParseInt(s, 10, 16)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetInt(int64(i))
			return nil
		}
	case *int32:
		return func(s string) error {
			i, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetInt(int64(i))
			return nil
		}
	case *int64:
		return func(s string) error {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetInt(int64(i))
			return nil
		}
	case *uint:
		return func(s string) error {
			i, err := strconv.ParseUint(s, 10, 0)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetUint(uint64(i))
			return nil
		}
	case *uint8:
		return func(s string) error {
			i, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetUint(uint64(i))
			return nil
		}
	case *uint16:
		return func(s string) error {
			i, err := strconv.ParseUint(s, 10, 16)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetUint(uint64(i))
			return nil
		}
	case *uint32:
		return func(s string) error {
			i, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetUint(uint64(i))
			return nil
		}
	case *uint64:
		return func(s string) error {
			i, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetUint(uint64(i))
			return nil
		}
	case *float32:
		return func(s string) error {
			i, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetFloat(float64(i))
			return nil
		}
	case *float64:
		return func(s string) error {
			i, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetFloat(float64(i))
			return nil
		}
	case *time.Duration:
		return func(s string) error {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			reflect.ValueOf(v).Elem().SetInt(int64(d))
			return nil
		}
	case *string:
		return func(s string) error {
			reflect.ValueOf(v).Elem().SetString(s)
			return nil
		}
	default:
		// TODO: handle other types
	}
	panic(fmt.Errorf("unsupported type %T", v))
}
