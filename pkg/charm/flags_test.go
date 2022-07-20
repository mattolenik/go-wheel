package charm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TODO: Add failure cases, wrong types, negatives into uints, etc

func TestVar(t *testing.T) {
	assert := assert.New(t)

	args := []string{
		"-string=two", "-bool=true", "-duration=5s", "-float64=3.14",
		"-int=-10", "-uint=10", "-int64=-30", "-uint64=30",
		"-intslice=1,2,3",
		"-int8slice=-4,5,6",
		"-int16slice=7,-8,9",
		"-int32slice=10,11,-12",
		"-int64slice=-13,-14,-15",
		"-uintslice=16,17,18",
		"-uint8slice=19,20,21",
		"-uint16slice=22, 23,24",
		"-uint32slice=25,  26,  27",
		"-uint64slice= 28,29,30 ",
		"-boolslice=true,false, true",
		"-stringslice=a,b,c",
		"-durationslice=1s,2m,3h",
		"-emptyslice=",
		"-blankslice=,,",
	}

	var i int
	var s string
	var b bool
	var d time.Duration
	var f float64
	var ui uint
	var i64 int64
	var ui64 uint64
	var emptySlice []int
	var blankSlice []string
	nopFn := func(c *Command) error { return nil }

	c := NewCommand("app", "app usage", nopFn)

	FlagVar(c, &s, "", false, "string", "")
	FlagVar(c, &b, false, false, "bool", "")
	FlagVar(c, &d, 0, false, "duration", "")
	FlagVar(c, &f, 0, false, "float64", "")
	FlagVar(c, &i, 0, false, "int", "")
	FlagVar(c, &ui, 0, false, "uint", "")
	FlagVar(c, &i64, 0, false, "int64", "")
	FlagVar(c, &ui64, 0, false, "uint64", "")

	FlagVar(c, &slices.Int, []int{}, false, "intslice", "")
	FlagVar(c, &slices.Int8, []int8{}, false, "int8slice", "")
	FlagVar(c, &slices.Int16, []int16{}, false, "int16slice", "")
	FlagVar(c, &slices.Int32, []int32{}, false, "int32slice", "")
	FlagVar(c, &slices.Int64, []int64{}, false, "int64slice", "")
	FlagVar(c, &slices.Uint, []uint{}, false, "uintslice", "")
	FlagVar(c, &slices.Uint8, []uint8{}, false, "uint8slice", "")
	FlagVar(c, &slices.Uint16, []uint16{}, false, "uint16slice", "")
	FlagVar(c, &slices.Uint32, []uint32{}, false, "uint32slice", "")
	FlagVar(c, &slices.Uint64, []uint64{}, false, "uint64slice", "")
	FlagVar(c, &slices.Bool, []bool{}, false, "boolslice", "")
	FlagVar(c, &slices.String, []string{}, false, "stringslice", "")
	FlagVar(c, &slices.Duration, []time.Duration{}, false, "durationslice", "")
	FlagVar(c, &emptySlice, []int{}, false, "emptyslice", "")
	FlagVar(c, &blankSlice, []string{"", "", ""}, false, "blankslice", "")

	err := c.Parse(args)
	assert.NoError(err)

	assert.Equal("two", s)
	assert.Equal(5*time.Second, d)
	assert.Equal(float64(3.14), f)
	assert.Equal(int(-10), i)
	assert.Equal(uint(10), ui)
	assert.Equal(int64(-30), i64)
	assert.Equal(uint64(30), ui64)

	assert.ElementsMatch([]int{1, 2, 3}, slices.Int)
	assert.ElementsMatch([]int8{-4, 5, 6}, slices.Int8)
	assert.ElementsMatch([]int16{7, -8, 9}, slices.Int16)
	assert.ElementsMatch([]int32{10, 11, -12}, slices.Int32)
	assert.ElementsMatch([]int64{-13, -14, -15}, slices.Int64)
	assert.ElementsMatch([]uint{16, 17, 18}, slices.Uint)
	assert.ElementsMatch([]uint8{19, 20, 21}, slices.Uint8)
	assert.ElementsMatch([]uint16{22, 23, 24}, slices.Uint16)
	assert.ElementsMatch([]uint32{25, 26, 27}, slices.Uint32)
	assert.ElementsMatch([]uint64{28, 29, 30}, slices.Uint64)
	assert.ElementsMatch([]bool{true, false, true}, slices.Bool)
	assert.ElementsMatch([]string{"a", "b", "c"}, slices.String)
	assert.ElementsMatch([]time.Duration{1 * time.Second, 2 * time.Minute, 3 * time.Hour}, slices.Duration)
}

var slices = struct {
	Int      []int
	Int8     []int8
	Int16    []int16
	Int32    []int32
	Int64    []int64
	Uint     []uint
	Uint8    []uint8
	Uint16   []uint16
	Uint32   []uint32
	Uint64   []uint64
	Bool     []bool
	Duration []time.Duration
	String   []string
}{
	Int:      []int{},
	Int8:     []int8{},
	Int16:    []int16{},
	Int32:    []int32{},
	Int64:    []int64{},
	Uint:     []uint{},
	Uint8:    []uint8{},
	Uint16:   []uint16{},
	Uint32:   []uint32{},
	Uint64:   []uint64{},
	Bool:     []bool{},
	Duration: []time.Duration{},
	String:   []string{},
}
