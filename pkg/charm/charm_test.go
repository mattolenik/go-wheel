package charm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	// c := "-a -b -c subcmd -d -e -f subcmd2 -g -h arg1 arg2"
	// cmd := Command{}

	var v int
	fd := &FlagDefinition[int]{
		Name:    "a",
		Usage:   "usage of a",
		Default: 5,
		Value:   &v,
	}
	fd.SetValue(8)
	assert.Equal(8, fd.GetValue())
	assert.Equal(8, *fd.Value)
}
