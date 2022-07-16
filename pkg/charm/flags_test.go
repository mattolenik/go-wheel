package charm

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVar(t *testing.T) {
	assert := assert.New(t)

	args := []string{
		"-string=two", "-bool=true", "-duration=5s", "-float64=3.14",
		"-int=-10", "-uint=10", "-int64=-30", "-uint64=30", "-intslice=4,2,5",
	}

	var i int
	var s string
	var b bool
	var d time.Duration
	var f float64
	var ui uint
	var i64 int64
	var ui64 uint64
	var sl []int

	c := NewCommand("app", "app usage")

	FlagVar(c, &s, "", "string", "")
	FlagVar(c, &b, false, "bool", "")
	FlagVar(c, &d, 0, "duration", "")
	FlagVar(c, &f, 0, "float64", "")
	FlagVar(c, &i, 0, "int", "")
	FlagVar(c, &ui, 0, "uint", "")
	FlagVar(c, &i64, 0, "int64", "")
	FlagVar(c, &ui64, 0, "uint64", "")
	FlagVar(c, &sl, []int{}, "intslice", "")

	err := c.Parse(args)
	assert.NoError(err)

	assert.Equal("two", s)
	assert.Equal(5*time.Second, d)
	assert.Equal(float64(3.14), f)
	assert.Equal(int(-10), i)
	assert.Equal(uint(10), ui)
	assert.Equal(int64(-30), i64)
	assert.Equal(uint64(30), ui64)
	assert.ElementsMatch([]int{4, 2, 5}, sl)

}

func TestVarFlagReturn(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-string=two", "-int=-10"}
	flags := flag.NewFlagSet("test", flag.PanicOnError)

	// TODO: seems like flags really shouldn't be coupled back to command? Find a better way to factor this
	c := NewCommand("app", "app usage")
	var i int
	var s string
	iFlag := FlagVar(c, &i, 1, "int", "")
	sFlag := FlagVar(c, &s, "str", "string", "")

	flags.Parse(args)

	assert.Equal("two", s)
	assert.Equal(int(-10), i)

	assert.Equal(i, *iFlag.Value, "Flag's value doesn't equal parsed value")
	assert.Equal(s, *sFlag.Value, "Flag's value doesn't equal parsed value")

	assert.Equal(&i, iFlag.Value, "Flag's value pointer doesn't point at the parsed result")
	assert.Equal(&s, sFlag.Value, "Flag's value pointer doesn't point at the parsed result")

	assert.Equal(1, iFlag.Default, "int default value did not match")
	assert.Equal("str", sFlag.Default, "string default value did not match")
}
