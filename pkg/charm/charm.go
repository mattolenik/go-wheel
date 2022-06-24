package charm

import (
	"fmt"

	"github.com/mattolenik/go-charm/internal/fun"
)

type Flag[T any] struct {
	Name         string
	Usage        string
	Value        *T
	DefaultValue T
}

type Command struct {
	Name     string
	Flags    []Flag[any]
	Commands []Command
}

func (c *Command) String() string {
	flagNames := fun.Map(c.Flags, func(f Flag[any]) string { return f.Name })
	subCommandNames := fun.Map(c.Commands, func(c Command) string { return c.Name })
	return fmt.Sprintf("Name: %q, Flags: %q, Subcommands: %q", c.Name, flagNames, subCommandNames)
}
