package charm

import (
	"flag"
	"fmt"

	"github.com/mattolenik/go-charm/internal/fn"
)

type Command struct {
	Name        string
	Usage       string
	Examples    []string
	Flags       []Flag
	Args        []string
	SubCommands []*Command
	flagSet     *flag.FlagSet
}

func (c *Command) String() string {
	subCommandNames := fn.Map(c.SubCommands, func(c *Command) string { return c.Name })
	return fmt.Sprintf("Name: %q, Flags: %q, Subcommands: %q", c.Name, c.Flags, subCommandNames)
}

func NewCommand(name, usage string) *Command {
	c := &Command{
		Name:  name,
		Usage: usage,
	}
	c.flagSet = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	return c
}

func (c *Command) SubCommand(name, usage string) *Command {
	subCommand := NewCommand(name, usage)
	c.SubCommands = append(c.SubCommands, subCommand)
	return subCommand
}

func (c *Command) Parse(args []string) error {
	fmt.Println(args)
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
	fmt.Println(args)
	return nil
}
