package main

import (
	"fmt"
	"os"

	"github.com/mattolenik/go-charm/pkg/charm"
)

func main() {
	err := mainE()
	if err == nil {
		return
	}
	panic(err)
}

func mainE() error {
	c := charm.NewCommand("dbdrawer", "dbd something something", func(c *charm.Command) error {
		fmt.Println("dbdrawer impl here")
		return nil
	})
	_ = c.SubCommand("sub1", "subcommand 1", func(c *charm.Command) error {
		fmt.Println("sub1 impl here")
		return nil
	})
	_ = c.SubCommand("sub2", "subcommand 2", func(c *charm.Command) error {
		fmt.Println("sub2 impl here")
		return nil
	})

	_ = charm.FlagF(c, 5, false, "intval", "the int value")
	err := c.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	return c.ExecDeepest()
}
