package wheel

import (
	"testing"
	"time"

	ppv3 "github.com/k0kubun/pp/v3"
	"github.com/mattolenik/go-charm/internal/fn"
	"github.com/stretchr/testify/assert"
)

var pp *ppv3.PrettyPrinter = ppv3.New()

func init() {
	pp.SetExportedOnly(true)
}

func TestParse(t *testing.T) {
	assert := assert.New(t)

	//args := []string{"-corge", "-grault", "-garply", "-waldo", "-fred", "-plugh", "-xyzzy", "-thud"}
	//args := []string{"-corge", "1", "-corge", "1b", "-waldo", "hidden", "-thud", "arg1", "arg2"}
	//args := []string{"-a", "1", "-b=4,5,6", "-c", "10", "-c", "-11", "-d", "20", "-d=30,40"}
	args := []string{"-a", "1", "-b=4,5,6", "-c", "10", "-c=-11", "-d", "20", "-d=30,40"}

	c := NewCommand("testparse", "testparse", "testparse", nil)

	aOpt := AddOption(c, false, fn.Ptr(5), "a", "a test")
	bOpt := AddOption(c, false, fn.Ptr([]int{}), "b", "b test")
	cOpt := AddOption(c, false, fn.Ptr([]int{}), "c", "c test")
	dOpt := AddOption(c, false, fn.Ptr([]int{}), "d", "d test")

	err := c.Parse(args)
	assert.NoError(err)

	assert.Equal(1, *aOpt.TypedValue)
	assert.Equal([]int{4, 5, 6}, *bOpt.TypedValue)
	assert.Equal([]int{10, -11}, *cOpt.TypedValue)
	assert.Equal([]int{20, 30, 40}, *dOpt.TypedValue)
}

// TODO: More thorough tests and parameterization.
func TestConvert(t *testing.T) {
	assert := assert.New(t)
	var x bool
	vf := converter(&x)
	assert.NoError(vf("false"))
	assert.Equal(false, x)
	assert.NoError(vf("1"))
	assert.Equal(true, x)

	j := JSON{}
	jj := converter(&j)
	assert.NoError(jj(`{"a":1}`))
	assert.Contains(j, "a")
	assert.Equal(float64(1), j["a"])

	{
		si := []bool{}
		c := converter(&si)
		assert.NoError(c("1, false,True"))
		assert.Equal([]bool{true, false, true}, si)
	}
	{
		si := []int{}
		c := converter(&si)
		assert.NoError(c("1,-2,3"))
		assert.Equal([]int{1, -2, 3}, si)
	}
	{
		si := []int8{}
		c := converter(&si)
		assert.NoError(c("1, -2, 3"))
		assert.Equal([]int8{1, -2, 3}, si)
	}
	{
		si := []int16{}
		c := converter(&si)
		assert.NoError(c("1,-2,3"))
		assert.Equal([]int16{1, -2, 3}, si)
	}
	{
		si := []int32{}
		c := converter(&si)
		assert.NoError(c("1,-2,3"))
		assert.Equal([]int32{1, -2, 3}, si)
	}
	{
		si := []int64{}
		c := converter(&si)
		assert.NoError(c("1,-2,3"))
		assert.Equal([]int64{1, -2, 3}, si)
	}

	{
		si := []uint{}
		c := converter(&si)
		assert.NoError(c("1,2,3"))
		assert.Equal([]uint{1, 2, 3}, si)
	}
	{
		si := []uint8{}
		c := converter(&si)
		assert.NoError(c("1,2,3"))
		assert.Equal([]uint8{1, 2, 3}, si)
	}
	{
		si := []uint16{}
		c := converter(&si)
		assert.NoError(c("1 , 2 , 3"))
		assert.Equal([]uint16{1, 2, 3}, si)
	}
	{
		si := []uint32{}
		c := converter(&si)
		assert.NoError(c("1,2,3"))
		assert.Equal([]uint32{1, 2, 3}, si)
	}
	{
		si := []uint64{}
		c := converter(&si)
		assert.NoError(c("1,2,3"))
		assert.Equal([]uint64{1, 2, 3}, si)
	}
	{
		si := []time.Duration{}
		c := converter(&si)
		assert.NoError(c("5s, 2m,3h"))
		assert.Equal([]time.Duration{5 * time.Second, 2 * time.Minute, 3 * time.Hour}, si)
	}
	{
		si := []string{}
		c := converter(&si)
		assert.NoError(c("a,b,c"))
		assert.Equal([]string{"a", "b", "c"}, si)
	}

	// Empty conditions
	{
		si := []string{}
		c := converter(&si)
		assert.NoError(c(""))
	}
	{
		si := []int{}
		c := converter(&si)
		assert.Error(c(""))
	}
}
