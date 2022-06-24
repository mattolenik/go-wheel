package charm

import (
	"flag"
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
	Usage    string
	Examples []string
	Flags    []string
	Commands []Command
	FlagSet  *flag.FlagSet
}

func (c *Command) String() string {
	subCommandNames := fun.Map(c.Commands, func(c Command) string { return c.Name })
	return fmt.Sprintf("Name: %q, Flags: %q, Subcommands: %q", c.Name, c.Flags, subCommandNames)
}

func NewCommand(flags *flag.FlagSet, name, usage string) *Command {
	return &Command{
		Name:     name,
		Usage:    usage,
		Examples: []string{},
		Commands: []Command{},
		FlagSet:  flags,
	}
}
