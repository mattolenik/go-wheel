package wheel

import (
	"strings"

	"github.com/mattolenik/go-charm/internal/fn"
)

func Parse(args []string) (flags fn.MultiMap[string, string], remainingArgs []string) {
	if len(args) == 0 {
		return
	}
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			remainingArgs = args[i:]
			return
		}
		var name string
		if strings.HasPrefix(arg, "--") {
			name = arg[2:]
		} else if strings.HasPrefix(arg, "-") {
			name = arg[1:]
		}
		nparts := strings.SplitN(name, "=", 1)
		if len(nparts) == 1 {
		}
	}
	return
}
