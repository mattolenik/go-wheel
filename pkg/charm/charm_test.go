package charm

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCharm(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-a=1", "-b=two", "-c=5s"}
	flags := flag.NewFlagSet("test", flag.PanicOnError)

	var a int
	var b string
	var c time.Duration

	Var(flags, &a, 0, "a", "a usage")
	Var(flags, &b, "", "b", "b usage")
	Var(flags, &c, 0, "c", "c usage")

	flags.Parse(args)

	assert.Equal(1, a)
	assert.Equal("two", b)
	assert.Equal(5*time.Second, c)
}
