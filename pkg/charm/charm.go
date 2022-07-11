package charm

import (
	"flag"
	"fmt"
	"reflect"
	"time"

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

type FlagType interface {
	string | float64 | bool | int | int64 | time.Duration
}

func FlagVar[T FlagType](c *Command, value *T, dfltValue T, name, usage string) *FlagDefinition[T] {
	f := &FlagDefinition[T]{
		Name:    name,
		Usage:   usage,
		Value:   value,
		Default: dfltValue,
	}
	if c.Flags == nil {
		c.Flags = []Flag{}
	}
	c.Flags = append(c.Flags, f)
	return f
}

type Flag interface {
	GetName() string
	GetUsage() string
	GetDefault() any
	GetValue() any
	SetValue(any)
}

type FlagDefinition[T FlagType] struct {
	Name    string
	Usage   string
	Default T
	Value   *T
}

func (fd *FlagDefinition[T]) GetName() string {
	return fd.Name
}

func (fd *FlagDefinition[T]) GetUsage() string {
	return fd.Name
}

func (fd *FlagDefinition[T]) GetDefault() any {
	return fd.Default
}

func (fd *FlagDefinition[T]) GetValue() any {
	return *fd.Value
}

func (fd *FlagDefinition[T]) SetValue(value any) {
	valType := reflect.TypeOf(fd.Value).Elem()
	vType := reflect.TypeOf(value)
	if vType.AssignableTo(valType) {
		v := reflect.ValueOf(fd.Value)
		v.Elem().Set(reflect.ValueOf(value))
	}
}

func FlagVars[T FlagType](c *Command, flags ...FlagDefinition[T]) {
	for _, flag := range flags {
		FlagVar(c, flag.Value, flag.Default, flag.Name, flag.Usage)
	}
}
