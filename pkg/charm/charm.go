package charm

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mattolenik/go-charm/internal/fn"
	"github.com/mattolenik/go-charm/internal/refract"
	"github.com/mattolenik/go-charm/pkg/twine"
	"github.com/mattolenik/go-charm/pkg/typ"
)

type CommandAction func(c *Command) error

type Command struct {
	Action         CommandAction
	Name           string
	Usage          string
	Examples       []string
	Args           []string
	SubCommands    []*Command
	Flags          []*Flag[any]
	root           *Command
	flagSet        *flag.FlagSet
	typeConverters map[reflect.Type]any
	Visited        bool
}

func NewCommand(name, usage string, action CommandAction) *Command {
	c := &Command{
		Action:         action,
		Name:           name,
		Usage:          usage,
		typeConverters: map[reflect.Type]any{},
		flagSet:        flag.NewFlagSet(name, flag.ContinueOnError),
	}
	c.root = c

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
	c.root.typeConverters[refract.TypeOf[T]()] = set
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

func (c *Command) SubCommand(name, usage string, action CommandAction) *Command {
	subCommand := &Command{
		Action:         action,
		Name:           name,
		Usage:          usage,
		typeConverters: map[reflect.Type]any{},
		root:           c.root,
		flagSet:        flag.NewFlagSet(c.Name, flag.ContinueOnError),
	}
	c.SubCommands = append(c.SubCommands, subCommand)
	return subCommand
}

func (c *Command) Parse(args []string) error {
	return c.parse(args)
}

func (c *Command) parse(args []string) error {
	c.Visited = true
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
		return subCmd.parse(args[1:])
	}

	// Otherwise, treat all the remaining args as arguments to the command.
	c.Args = args
	return nil
}

// TreeString creates an indented, human readable multi-line string representation of the command tree.
func (c *Command) TreeString(indent string) string {
	sb := &strings.Builder{}
	c.Walk(func(depth int, c *Command) {
		sb.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(indent, depth), c.String()))
	})
	return sb.String()
}

func (c *Command) Walk(fn func(depth int, c *Command)) {
	c.walk(0, fn)
}

func (c *Command) walk(depth int, fn func(int, *Command)) {
	fn(depth, c)
	for _, subCmd := range c.SubCommands {
		fn(depth+1, subCmd)
	}
}

func (c *Command) WalkVisited(fn func(depth int, c *Command)) {
	c.walkVisited(0, fn)
}

func (c *Command) walkVisited(depth int, fn func(int, *Command)) {
	fn(depth, c)
	for _, subCmd := range c.SubCommands {
		if subCmd.Visited {
			fn(depth+1, subCmd)
		}
	}
}

func (c *Command) ChosenCommand() *Command {
	chosen := c
	deepest := -1
	c.WalkVisited(func(depth int, c *Command) {
		if depth > deepest {
			deepest = depth
			chosen = c
		}
	})
	return chosen
}

func (c *Command) Exec() error {
	if c.Action == nil {
		return fmt.Errorf("command %q has no Action set, cannot execute it", c.Name)
	}
	return c.Action(c)
}

func (c *Command) Deepest(fn func(*Command) bool) *Command {
	d := Command{}
	deepest := 0
	c.deepest(&d, &deepest, 0, fn)
	return &d
}

func (c *Command) deepest(deepestMatch *Command, deepest *int, depth int, fn func(*Command) bool) {
	if fn(c) {
		*deepestMatch = *c
		*deepest = depth
	}
	for _, sc := range c.SubCommands {
		sc.deepest(deepestMatch, deepest, depth+1, fn)
	}
}

// TODO: Rename this
func (c *Command) ExecDeepest() error {
	deepest := c.Deepest(func(c *Command) bool {
		return c.Visited
	})
	return deepest.Exec()
}
