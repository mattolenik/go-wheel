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
	sc := c.SubCommand("sub1", "subcommand 1")

	_ = charm.FlagF(c, 5, "intval", "the int value")
	err := c.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	c.Walk(0, func(d int, c *charm.Command) {
		if c == sc && c.Visited {
			fmt.Println(c)
		}
	})
	return nil
}
