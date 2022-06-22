package charm

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCharm(t *testing.T) {
	assert := assert.New(t)

	args := []string{"-a=1", "-b=two", "-c=5s", "-d=3.14"}
	flags := flag.NewFlagSet("test", flag.PanicOnError)

	var a int
	var b string
	var c time.Duration
	var d float64

	Var(flags, &a, 0, "a", "a usage")
	Var(flags, &b, "", "b", "b usage")
	Var(flags, &c, 0, "c", "c usage")
	Var(flags, &d, 0, "d", "d usage")

	flags.Parse(args)

	assert.Equal(1, a)
	assert.Equal("two", b)
	assert.Equal(5*time.Second, c)
	assert.Equal(float64(3.14), d)
}
