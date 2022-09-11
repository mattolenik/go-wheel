package main

import (
	"fmt"
	"os"

	ppv3 "github.com/k0kubun/pp/v3"
	"github.com/mattolenik/go-wheel/pkg/wheel"
)

var pp *ppv3.PrettyPrinter = ppv3.New()

func init() {
	pp.SetExportedOnly(true)
}

func main() {
	err := mainE()
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
}

func mainE() error {
	args := []string{"-abc=123", "-def", "-y", "456"}
	root := wheel.NewCommand("myapp", "myapp <args>", "does appy things", nil)
	_ = root.SubCommand("sub", "sub <args>", "does subby things", nil)
	var abc int
	var y int
	var def bool
	abcOpt := wheel.AddOption[int](root, "abc", "does abc things").Bind(&abc)
	wheel.AddOption[int](root, "y", "does y things").Bind(&y)
	wheel.AddOption[bool](root, "def", "def on off").Bind(&def)
	err := root.Parse(args)
	if err != nil {
		return err
	}
	fmt.Printf("abc: %d, y: %d, def: %v\n", abc, y, def)
	pp.Println(abcOpt)
	return nil
}
