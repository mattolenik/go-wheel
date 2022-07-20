package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/mattolenik/go-charm/pkg/charm"
)

func main() {
	err := mainE()
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
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
	err := StructCmd(c, dbf)
	if err != nil {
		return err
	}

	err = c.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	return nil
}

type DbdrawerFlags struct {
	Sl []int  `flag:"sl" desc:"Int slicerydicer" usage:"some information usage might go here"`
	A  string `flag:"a" required:"true"`
}

type FlagTags struct {
	Flag, Desc, Usage string
	Required          bool
}

// String implements the fmt.Stringer interface for FlagTags
func (ft *FlagTags) String() string {
	return fmt.Sprintf("Flag: %q, Desc: %q, Usage: %q, Required: %v", ft.Flag, ft.Desc, ft.Usage, ft.Required)
}

// TODO: this whole thing is mistakenly relating the flag tags to entire struct and not each field, fix and clarify this
// Parse takes all the tags from all the fields of the given struct and assigns them to the corresponding fields of a FlagTags struct.
func Parse(struc any) (*FlagTags, error) {
	strucType := reflect.TypeOf(struc).Elem()
	if strucType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("func FlagTags.Parse expected a value of kind struct, instead got %s", strucType.Kind())
	}
	for i := 0; i < strucType.NumField(); i++ {
		f := strucType.Field(i)
		tag := f.Tag
		reqTagName := "required"
		tagVal := tag.Get(reqTagName)
		requiredBool, err := strconv.ParseBool(tagVal)
		if err != nil {
			return nil, fmt.Errorf("invalid value for boolean: %s", "required")
		}
		return &FlagTags{
			Flag:     tag.Get("flag"),
			Desc:     tag.Get("desc"),
			Usage:    tag.Get("usage"),
			Required: requiredBool,
		}, nil
	}
	return nil
}

func StructCmd(c *charm.Command, struc any) error {
	strucType := reflect.TypeOf(struc).Elem()
	if strucType.Kind() != reflect.Struct {
		return fmt.Errorf("func StructCmd expected a value of kind struct, instead got %s", strucType.Kind())
	}
	for i := 0; i < strucType.NumField(); i++ {
		f := strucType.Field(i)
		tag := f.Tag
		flag, desc, usage, required := tag.Get("flag"), tag.Get("desc"), tag.Get("usage"), tag.Get("required")
	}
	return nil
}
