package charm

import (
	"flag"
	"fmt"

	"github.com/mattolenik/go-charm/internal/fn"
)

type Command struct {
	Name     string
	Usage    string
	Examples []string
	Flags    []Flag
	Args     []string
	Commands []Command
	FlagSet  *flag.FlagSet
}

func (c *Command) String() string {
	subCommandNames := fn.Map(c.Commands, func(c Command) string { return c.Name })
	return fmt.Sprintf("Name: %q, Flags: %q, Subcommands: %q", c.Name, c.Flags, subCommandNames)
}

func (c *Command) Parse(args []string) error {
	// c.FlagSet = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	// for _, f := range c.Flags {
	// 	//c.FlagSet.String(f, "", "")
	// }
	// c.FlagSet.String()
	// err := c.FlagSet.Parse(args)
	// if err != nil {
	// 	return err
	// }
	return nil
}
