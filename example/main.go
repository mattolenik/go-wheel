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
	c := charm.NewCommand("dbdrawer", "dbd something something")
	_ = c.SubCommand("sub1", "subcommand 1")
	_ = c.SubCommand("sub2", "subcommand 2")

	_ = charm.FlagF(c, 5, "intval", "the int value")
	err := c.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	c.WalkVisited(0, func(d int, c *charm.Command) {
		fmt.Println(c)
	})
	fmt.Println()
	fmt.Println(c.ChosenCommand().String())
	return nil
}
