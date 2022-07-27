package wheel

import (
	"testing"

	ppv3 "github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
)

var pp *ppv3.PrettyPrinter = ppv3.New()

func TestParse(t *testing.T) {
	assert := assert.New(t)

	//args := []string{"-corge", "-grault", "-garply", "-waldo", "-fred", "-plugh", "-xyzzy", "-thud"}
	args := []string{"-corge", "1", "-corge", "1b", "-waldo", "hidden", "-thud", "arg1", "arg2"}

	c := NewCommand("testparse", "testparse", "testparse", nil)
	err := c.Parse(args)
	assert.NoError(err)
	pp.Println(c)
}
