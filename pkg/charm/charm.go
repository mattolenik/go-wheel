package charm

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mattolenik/go-charm/internal/fn"
	"github.com/mattolenik/go-charm/internal/twine"
	"github.com/mattolenik/go-charm/internal/typ"
)

type Command struct {
	Name           string
	Usage          string
	Examples       []string
	Args           []string
	SubCommands    []*Command
	Flags          []*Flag[any]
	root           *Command
	flagSet        *flag.FlagSet
	typeConverters map[reflect.Type]any
}

func NewCommand(name, usage string) *Command {
	c := &Command{
		Name:           name,
		Usage:          usage,
		typeConverters: map[reflect.Type]any{},
	}
	c.root = c
	c.flagSet = flag.NewFlagSet(c.Name, flag.ContinueOnError)

	// TODO: Get these implemented natively rather than registered as type converters.
	RegisterTypeConverter(c, listConverter[int])
	RegisterTypeConverter(c, listConverter[int8])
	RegisterTypeConverter(c, listConverter[int16])
	RegisterTypeConverter(c, listConverter[int32])
	RegisterTypeConverter(c, listConverter[int64])
	RegisterTypeConverter(c, listConverter[uint])
	RegisterTypeConverter(c, listConverter[uint8])
	RegisterTypeConverter(c, listConverter[uint16])
	RegisterTypeConverter(c, listConverter[uint32])
	RegisterTypeConverter(c, listConverter[uint64])
	RegisterTypeConverter(c, listConverter[bool])
	RegisterTypeConverter(c, listConverter[time.Duration])
	RegisterTypeConverter(c, listConverter[string])
	return c
}

func listConverter[T typ.StringRepresentable](s string, a *[]T) error {
	var err error
	*a, err = twine.FromDelimetedList[T](s, ",")
	return err
}

// FlagSetterFunc describes a function that takes a string value from a flag at the command line.
// Given -foo=bar, this function will be passed the string "bar" and a pointer to where the parsed value should be placed.
type FlagSetterFunc[T any] func(stringValue string, value *T) error

func RegisterTypeConverter[T any](c *Command, set FlagSetterFunc[T]) {
	c.root.typeConverters[fn.TypeOf[T]()] = set
}

func (c *Command) FindTypeConverter(t reflect.Type) (any, bool) {
	cv, ok := c.root.typeConverters[t]
	if !ok {
		return nil, false
	}
	return cv, true
}

func (c *Command) String() string {
	return fmt.Sprintf("Name: %q, Flags: %q, Subcommands: %d", c.Name, c.Flags, len(c.SubCommands))
}

func (c *Command) SubCommand(name, usage string) *Command {
	subCommand := NewCommand(name, usage)
	subCommand.root = c.root
	c.SubCommands = append(c.SubCommands, subCommand)
	return subCommand
}

func (c *Command) Parse(args []string) error {
	if len(args) == 0 {
		return nil
	}
	err := c.flagSet.Parse(args)
	if err == flag.ErrHelp {
		return fmt.Errorf("command %q: %w", c.Name, err)
	}
	if err != nil {
		return fmt.Errorf("error parsing flags for command %q: %w", c.Name, err)
	}
	c.flagSet.NFlag()
	remaining := c.flagSet.NArg()
	args = args[len(args)-remaining:]
	if len(args) == 0 {
		return nil
	}

	// Check whether the first arg is the name of a subcommand. If so, treat it as an invocation of that subcommand.
	// All arguments from here on out are handled by the subcommand.
	subCmdName := args[0]
	if ok, subCmd := fn.Find(c.SubCommands, func(sc *Command) bool { return sc.Name == subCmdName }); ok {
		return subCmd.Parse(args[1:])
	}

	// Otherwise, treat all the remaining args as arguments to the command.
	c.Args = args
	return nil
}

func (c *Command) TreePrint(indent string) string {
	sb := &strings.Builder{}
	c.walk(0, func(depth int, c *Command) {
		sb.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(indent, depth), c.String()))
	})
	return sb.String()
}

func (c *Command) walk(depth int, fn func(int, *Command)) {
	fn(depth, c)
	for _, subCmd := range c.SubCommands {
		fn(depth+1, subCmd)
	}
}
