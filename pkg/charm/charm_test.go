package charm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	//args := []string{"-a", "1", "subcmd"}
	args := []string{"-a", "1,2,3"}
	c := NewCommand("app", "app usage")

	var v int
	FlagVar(c, &v, 5, "a", "usage of a")

	assert.NoError(c.Parse(args))
}
