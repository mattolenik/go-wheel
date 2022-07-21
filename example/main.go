package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/k0kubun/pp/v3"
	"github.com/mattolenik/go-charm/pkg/charm"
	"github.com/mattolenik/go-charm/pkg/typ"
)

func main() {
	err := mainE()
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
}

func mainE() error {
	pp := pp.New()
	pp.SetExportedOnly(true)

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

	//_ = charm.FlagF(c, []int{}, false, "sl", "a slice")
	//_ = charm.FlagF(c, "", false, "a", "a string")

	//err := c.Parse(os.Args[1:])
	//if err != nil {
	//	return err
	//}

	//_ = charm.FlagF(c, 5, false, "intval", "the int value")
	//err = c.Parse(os.Args[1:])
	//if err != nil {
	//	return err
	//}
	//return c.ExecDeepest()

	dbf := &DbdrawerFlags{
		Sl: []int{5, 9},
		A:  "",
	}
	ft, err := Parse(dbf)
	if err != nil {
		return err
	}

	err = c.Parse(os.Args[1:])
	if err != nil {
		return err
	}
	pp.Println(ft)
	fmt.Println("----------------------------------------------------")
	pp.Println(c)
	return nil
}

type DbdrawerFlags struct {
	Sl []int  `flag:"sl" desc:"Int slicerydicer" usage:"some information usage might go here"`
	A  string `flag:"a" required:"true"`
}

type FlagTags struct {
	Flag, Desc, Usage *string
	Required          bool
}

// String implements the fmt.Stringer interface for FlagTags
func (ft *FlagTags) String() string {
	return fmt.Sprintf("Flag: %q, Desc: %q, Usage: %q, Required: %v", *ft.Flag, *ft.Desc, *ft.Usage, ft.Required)
}

// TODO: this whole thing is mistakenly relating the flag tags to entire struct and not each field, fix and clarify this
// Parse takes all the tags from all the fields of the given struct and assigns them to the corresponding fields of a FlagTags struct.
func Parse(c *charm.Command, struc any) (map[*reflect.StructField]FlagTags, error) {
	strucType := reflect.TypeOf(struc).Elem()
	if strucType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("func FlagTags.Parse expected a value of kind struct, instead got %s", strucType.Kind())
	}
	result := map[*reflect.StructField]FlagTags{}
	for i := 0; i < strucType.NumField(); i++ {
		f := strucType.Field(i)

		required := false
		requiredTag, ok := f.Tag.Lookup("required")
		if ok {
			var err error
			required, err = typ.Parse[bool](requiredTag)
			if err != nil {
				return nil, fmt.Errorf("value %q is not a bool", requiredTag)
			}
		}
		result[&f] = FlagTags{
			Flag:     Ptr(f.Tag.Get("flag")),
			Desc:     Ptr(f.Tag.Get("desc")),
			Usage:    Ptr(f.Tag.Get("usage")),
			Required: required,
		}
	}
	return result, nil
}

func Ptr[T any](v T) *T {
	return &v
}
