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
		"-int=-10", "-uint=10", "-int64=-30", "-uint64=30",
	}
	flags := flag.NewFlagSet("test", flag.PanicOnError)

	var i int
	var s string
	var b bool
	var d time.Duration
	var f float64
	var ui uint
	var i64 int64
	var ui64 uint64

	Var(flags, &s, "", "string", "")
	Var(flags, &b, false, "bool", "")
	Var(flags, &d, 0, "duration", "")
	Var(flags, &f, 0, "float64", "")
	Var(flags, &i, 0, "int", "")
	Var(flags, &ui, 0, "uint", "")
	Var(flags, &i64, 0, "int64", "")
	Var(flags, &ui64, 0, "uint64", "")

	flags.Parse(args)

	assert.Equal("two", s)
	assert.Equal(5*time.Second, d)
	assert.Equal(float64(3.14), f)
	assert.Equal(int(-10), i)
	assert.Equal(uint(10), ui)
	assert.Equal(int64(-30), i64)
	assert.Equal(uint64(30), ui64)
}

func TestVarFlagReturn(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-string=two", "-int=-10"}
	flags := flag.NewFlagSet("test", flag.PanicOnError)

	var i int
	var s string
	iFlag := Var(flags, &i, 1, "int", "")
	sFlag := Var(flags, &s, "str", "string", "")

	flags.Parse(args)

	assert.Equal("two", s)
	assert.Equal(int(-10), i)

	assert.Equal(i, *iFlag.Value, "Flag's value doesn't equal parsed value")
	assert.Equal(s, *sFlag.Value, "Flag's value doesn't equal parsed value")

	assert.Equal(&i, iFlag.Value, "Flag's value pointer doesn't point at the parsed result")
	assert.Equal(&s, sFlag.Value, "Flag's value pointer doesn't point at the parsed result")

	assert.Equal(1, iFlag.DefaultValue, "int default value did not match")
	assert.Equal("str", sFlag.DefaultValue, "string default value did not match")
}
