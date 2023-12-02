package wheel

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mattolenik/go-wheel/internal/fn"
	"github.com/mattolenik/go-wheel/internal/refract"
)

type Command struct {
	Name        string
	Description string
	Usage       string
	Examples    []string
	Options     []Option
	SubCommands []*Command
	// TODO: make arg parsing strongly typed
	Args            []string
	parent          *Command
	Invoked         bool
	SeparateRawArgs bool
	RawArgSeparator string
}

func NewCommand(name, usage, description string, examples []string) *Command {
	c := &Command{
		Name:            name,
		Usage:           usage,
		Description:     description,
		Examples:        examples,
		SeparateRawArgs: true,
		RawArgSeparator: "--",
	}
	return c
}

func AddOption[T CommandLineType](c *Command, name, description string) *TypedOption[T] {
	if _, exists := fn.FindP(c.Options, func(o *Option) bool { return o.Name == name }); exists {
		panic(fmt.Errorf("an option with the name %q already exists", name))
	}
	var value T
	o := &TypedOption[T]{Value: &value}
	o.Option = Option{
		Name:        name,
		Description: description,
		Type:        refract.TypeOf[T](),
		// TODO: converter needs to reference this option instance so that it can add to it, e.g. -a=1 -a=2 should be an array
		Setter: converter(&value),
		// TODO: revisit this odd pattern
		Get:         func() any { return value },
		TypedOption: o,
	}
	c.Options = append(c.Options, o.Option)
	return o
}

func (c *Command) SubCommand(name, usage, description string, examples []string) *Command {
	if _, ok := fn.Find(c.SubCommands, func(sc *Command) bool { return sc.Name == name }); ok {
		panic(fmt.Errorf("a subcommand with the name %q already exists", name))
	}
	sc := NewCommand(name, usage, description, examples)
	sc.parent = c
	c.SubCommands = append(c.SubCommands, sc)
	return sc
}

func (c *Command) gatherGlobalOpts() []*Option {
	opts := fn.Filter(c.Options, func(o *Option) bool { return o.Global })
	parentOpts := []*Option{}
	if c.parent != nil {
		parentOpts = c.parent.gatherGlobalOpts()
	}
	for _, o := range parentOpts {
		_, found := fn.Find(opts, func(oo *Option) bool { return oo.Name == o.Name })
		if !found {
			opts = append(opts, o)
		}
	}
	return opts
}

func (c *Command) parseBetter(args []string) error {
	if len(args) == 0 {
		return nil
	}
	arg := args[0]
	if c.SeparateRawArgs && arg == c.RawArgSeparator {
		c.Args = append(c.Args, args[1:]...)
		return nil
	}
	if arg == c.RawArgSeparator || arg == "-" {
		c.Args = append(c.Args, arg)
		return c.parseBetter(args[1:])
	}
	if sc, ok := fn.Find(c.SubCommands, func(c *Command) bool { return c.Name == arg }); ok {
		return sc.parseBetter(args[1:])
	}
	if trimmed := strings.TrimPrefix(arg, "--"); trimmed != arg {
		arg = trimmed
	} else if trimmed := strings.TrimPrefix(arg, "-"); trimmed != arg {
		arg = trimmed
	} else {
		c.Args = append(c.Args, arg)
		return c.parseBetter(args[1:])
	}

	opt, ok := fn.FindP(c.Options, func(o *Option) bool { return o.Name == arg })
	if !ok {
		// TODO: Also search upwards for global options
		return &InvalidOptionError{arg}
	}
	val, hasVal, remainingArgs := parseOptValue(arg, args[1:])
	if hasVal {
		err := opt.Setter(val)
		if err != nil {
			return fmt.Errorf("invalid value for %q, %w", arg, err)
		}
	} else {
		// If there is no value, just the bare flag, this is treated as a boolean type
		err := opt.Setter("True")
		if err != nil {
			return fmt.Errorf("invalid value for %q, %w", arg, err)
		}
	}
	return c.parseBetter(remainingArgs)
}

func parseOptValue(opt string, nextArgs []string) (string, bool, []string) {
	parts := strings.SplitN(opt, "=", 2)
	if len(parts) == 2 {
		return parts[1], true, nextArgs
	}
	if len(nextArgs) == 0 {
		return "", false, nextArgs
	}
	if strings.HasPrefix(nextArgs[0], "-") {
		return "", false, nextArgs
	}
	return nextArgs[0], true, nextArgs[1:]
}

func (c *Command) Parse(args []string) error {
	opts, remaining := parseOptions(args)
	fmt.Println(opts)
	for _, definedOpt := range c.Options {
		opt := definedOpt.Name
		fmt.Println("opt: " + opt)
		values, ok := opts.Lookup(opt)
		if !ok && definedOpt.IsRequired {
			if definedOpt.Default != nil {
				if err := definedOpt.Setter(*definedOpt.Default); err != nil {
					return fmt.Errorf("option %q: %w", opt, err)
				}
				continue
			}
			return fmt.Errorf("must provide a value for option %q", opt)
		}
		supportedOpts := fn.Filter(c.Options, func(o *Option) bool { return (*o).Name == opt })
		if len(supportedOpts) == 0 {
			return fmt.Errorf("unsupported option %q, did you mean <TODO: insert help here>?", opt)
		}
		if len(supportedOpts) > 1 {
			// This is a panic rather than an error because duplicate options indicate a serious bug in the program.
			panic(fmt.Errorf("duplicate option found, %q was defined %d times, must be only once", opt, len(supportedOpts)))
		}
		o := *supportedOpts[0]
		hasValue := fn.Has(values.Values(), func(v string) bool { return v != "" })
		fmt.Println(hasValue)
		if !hasValue && o.Default != nil {
			if err := o.Setter(*o.Default); err != nil {
				return fmt.Errorf("option %q: %w", opt, err)
			}
			continue
		}
		if o.Type.Kind() == reflect.Slice {
			if len(values) == 0 {
				return fmt.Errorf("option %q requires at least one value", opt)
			}
			for v := range values {
				if err := o.Setter(v); err != nil {
					return fmt.Errorf("option %q: %w", opt, err)
				}
			}
			continue
		}
		if len(values) == 0 {
			if o.Type == refract.TypeOf[bool]() {
				if err := o.Setter("true"); err != nil {
					return fmt.Errorf("option %q: %w", opt, err)
				}
				continue
			}
			return fmt.Errorf("option %q requires a value", opt)
		}
		if len(values) == 1 {
			v := values.Values()[0]
			if err := o.Setter(v); err != nil {
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

// parseOptions takes CLI arguments and returns a mapping of options to values, plus the remaining, non-option arguments.
// Example:
//
//	[]string{"-a=1", "-b=2", "-c=3", "-b=4", "arg1", "arg2", "arg3"}
//	  becomes:
//	map[string]Set[string]{"a":{"1"}, "b":{"2", "4"}, "c":{"3"}}, []string{"arg1", "arg2", "arg3"}
func parseOptions(args []string) (fn.MultiMap[string, string], []string) {
	if len(args) == 0 {
		return fn.MultiMap[string, string]{}, args
	}
	if s := strings.SplitN("a=b", "=", 1); len(s) > 1 {
	}
	flags := fn.MultiMap[string, string]{}
	var i int
	for i := range args {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			arg = arg[2:]
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
		}
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 1 { // This is the case of an = sign NOT appearing in the argument, i.e. NOT -a=b

			// Look ahead at the next value to see if it is a flag or not.
			// If not, treat it as an argument for the current flag.
			// If there is no next value then this is the end of the list and we can just continue out of here.
			flag := parts[0]
			value, ok := fn.Index(args, i+1)
			if !ok {
				// Continue if this is the end of the args list
				continue
			}
			if strings.HasPrefix(value, "-") {
				// Next arg is flag
				continue
			}
			// Next arg is the flag's value
			flags.Put(flag, value)
		} else if len(parts) == 2 { // This is the case of an = sign being present, i.e. -a=b
			// The value on the right side of the = sign is the flag's value
			flags.Put(parts[0], parts[1])
		} else {
			// This shouldn't be possible since strings.SplitN(2) should never return a slice of length > 2
			panic(fmt.Errorf("expected only 2 strings to be returned by strings.SplitN but instead got %d", len(parts)))
		}
	}
	return flags, args[i:]
}

func parseAndAppend[T CommandLineType](sep string, v *[]T) func(string) error {
	return func(s string) error {
		parts := strings.Split(strings.TrimSpace(s), sep)
		if len(*v) == 0 {
			*v = make([]T, 0, len(parts))
		}
		for _, part := range parts {
			part = strings.TrimSpace(part) // allows the user to use spaces between separator, e.g. "1, 2, 3"
			var value T
			f := converter(&value)
			if err := f(part); err != nil {
				return fmt.Errorf("error decoding %q: %w", s, err)
			}
			*v = append(*v, value)
		}
		return nil
	}
}

func converter[T CommandLineType](v *T) func(string) error {
	switch v := any(v).(type) {
	case *[]bool:
		return parseAndAppend(",", v)
	case *[]int:
		return parseAndAppend(",", v)
	case *[]int8:
		return parseAndAppend(",", v)
	case *[]int16:
		return parseAndAppend(",", v)
	case *[]int32:
		return parseAndAppend(",", v)
	case *[]int64:
		return parseAndAppend(",", v)
	case *[]uint:
		return parseAndAppend(",", v)
	case *[]uint8:
		return parseAndAppend(",", v)
	case *[]uint16:
		return parseAndAppend(",", v)
	case *[]uint32:
		return parseAndAppend(",", v)
	case *[]uint64:
		return parseAndAppend(",", v)
	case *[]time.Duration:
		return parseAndAppend(",", v)
	case *[]string:
		return parseAndAppend(",", v)
	case *[]any:
		return parseAndAppend(",", v)

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
	case *JSON:
		return func(s string) error {
			return v.FromString(s)
		}
	default:
		// TODO: handle other types
	}
	panic(fmt.Errorf("unsupported type %T", v))
}
