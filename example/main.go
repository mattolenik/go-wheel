package main

import (
	"fmt"
	"os"
	"reflect"

	ppv3 "github.com/k0kubun/pp/v3"
	"github.com/mattolenik/go-wheel/pkg/typ"
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

// Flag is bool and others are options?
type Flag struct {
}

func mainE() error {
	args := []string{"-abc=123", "-def", "-x", "-abc=5", "-y", "456"}
	// Pass in multimap instead of receiving, pass in expected flags and flag info?
	flags, remainingArgs := wheel.ParseFlags(args)
	fmt.Println("Flags:")
	pp.Println(flags)
	fmt.Println()
	fmt.Println("Remaining args:")
	pp.Println(remainingArgs)
	return nil
}

func Index[T any](slice []T, i int) (v T, ok bool) {
	if i >= len(slice) {
		return
	}
	return slice[i], true
}

// func mainE() error {
// 	pp := pp.New()
// 	pp.SetExportedOnly(true)

// 	c := wheel.NewCommand("dbdrawer", "dbd something something", func(c *wheel.Command) error {
// 		fmt.Println("dbdrawer impl here")
// 		return nil
// 	})
// 	_ = c.SubCommand("sub1", "subcommand 1", func(c *wheel.Command) error {
// 		fmt.Println("sub1 impl here")
// 		return nil
// 	})
// 	_ = c.SubCommand("sub2", "subcommand 2", func(c *wheel.Command) error {
// 		fmt.Println("sub2 impl here")
// 		return nil
// 	})

// 	//_ = wheel.FlagF(c, []int{}, false, "sl", "a slice")
// 	//_ = wheel.FlagF(c, "", false, "a", "a string")

// 	//err := c.Parse(os.Args[1:])
// 	//if err != nil {
// 	//	return err
// 	//}

// 	//_ = wheel.FlagF(c, 5, false, "intval", "the int value")
// 	//err = c.Parse(os.Args[1:])
// 	//if err != nil {
// 	//	return err
// 	//}
// 	//return c.ExecDeepest()

// 	dbf := &DbdrawerFlags{
// 		Sl: []int{5, 9},
// 		A:  "",
// 	}
// 	err := Parse(c, dbf)
// 	if err != nil {
// 		return err
// 	}
// 	pp.Println(dbf)
// 	err = c.Parse(os.Args[1:])
// 	if err != nil {
// 		return err
// 	}
// 	// pp.Println(ft)
// 	// fmt.Println("----------------------------------------------------")
// 	// pp.Println(c)
// 	return nil
// }

type DbdrawerFlags struct {
	Sl []int  `flag:"sl" desc:"Int slicerydicer" usage:"some information usage might go here"`
	A  string `flag:"a" desc:"test flag" usage:"usage here" required:"true"`
}

type FlagTags struct {
	Flag, Desc, Usage *string
	Required          bool
}

// String implements the fmt.Stringer interface for FlagTags
func (ft *FlagTags) String() string {
	return fmt.Sprintf("Flag: %q, Desc: %q, Usage: %q, Required: %v", *ft.Flag, *ft.Desc, *ft.Usage, ft.Required)
}

// p := reflect.New(reflect.TypeOf(v))
// p.Elem().Set(reflect.ValueOf(v))

// return p.Interface()

// Parse takes all the tags from all the fields of the given struct and assigns them to the corresponding fields of a FlagTags struct.
func Parse(c *wheel.Command, struc any) error {
	strucType := reflect.TypeOf(struc).Elem()
	structVal := reflect.ValueOf(struc).Elem()
	if strucType.Kind() != reflect.Struct {
		return fmt.Errorf("func FlagTags.Parse expected a value of kind struct, instead got %s", strucType.Kind())
	}
	for i := 0; i < strucType.NumField(); i++ {
		var aVal *any
		f := strucType.Field(i)
		sv := structVal.Field(i)
		svi := sv.Interface()
		aVal = &svi

		required := false
		requiredTag, ok := f.Tag.Lookup("required")
		if ok {
			var err error
			required, err = typ.Parse[bool](requiredTag)
			if err != nil {
				return fmt.Errorf("value %q is not a bool", requiredTag)
			}
		}
		ftgs := FlagTags{
			Flag:     Ptr(f.Tag.Get("flag")),
			Desc:     Ptr(f.Tag.Get("desc")),
			Usage:    Ptr(f.Tag.Get("usage")),
			Required: required,
		}
		wheel.FlagVar(c, aVal, nil, ftgs.Required, *ftgs.Flag, *ftgs.Usage)
	}
	return nil
}

func Ptr[T any](v T) *T {
	return &v
}
