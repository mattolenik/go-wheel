package charm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	//args := []string{"-a", "1", "subcmd"}
	args := []string{"-a", "subcmd"}
	c := NewCommand("app", "app usage")

	var v int
	// fd := &FlagDefinition[int]{
	// 	Name:    "a",
	// 	Usage:   "usage of a",
	// 	Default: 5,
	// 	Value:   &v,
	// }
	// fd.SetValue(8)
	// assert.Equal(8, fd.GetValue())
	// assert.Equal(8, *fd.Value)
	FlagVar(c, &v, 5, "a", "usage of a")

	assert.NoError(c.Parse(args))
}
