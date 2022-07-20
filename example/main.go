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

	_ = charm.FlagF(c, 5, "intval", "the int value")
	err := c.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	fmt.Println(c.TreePrint("  "))
	return nil
}
