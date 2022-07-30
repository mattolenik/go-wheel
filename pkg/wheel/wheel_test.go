package wheel

import (
	"testing"

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
	args := []string{"-a", "1"}

	c := NewCommand("testparse", "testparse", "testparse", nil)
	aOpt := AddOption(c, false, fn.Ptr(5), "a", "a test")
	err := c.Parse(args)
	assert.NoError(err)
	assert.Equal(1, *aOpt.TypedValue)
	pp.Println(c)
}

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

	si := []int{}
	sii := converter(&si)
	err := sii("1,2,3")
	assert.NoError(err)
	assert.Equal([]int{1, 2, 3}, si)
}
