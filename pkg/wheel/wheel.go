package wheel

import (
	"fmt"
	"strings"

	"github.com/mattolenik/go-charm/internal/fn"
)

func ParseFlags(args []string) (fn.MultiMap[string, string], []string) {
	if len(args) == 0 {
		return fn.MultiMap[string, string]{}, args
	}
	flags := fn.MultiMap[string, string]{}
	var i int
	var arg string
	for i, arg = range args {
		if strings.HasPrefix(arg, "--") {
			arg = arg[2:]
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
		} else {
			return flags, args[i:]
		}
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 1 {
			flag := parts[0]
			value, ok := Index(args, i+1)
			if !ok {
				// Continue if this is the end of the list
				continue
			}
			if strings.HasPrefix(value, "-") {
				// Next arg is flag
				continue
			}
			flags.Put(flag, value)
		} else if len(parts) == 2 {
			flags.Put(parts[0], parts[1])
		} else {
			// This shouldn't be possible since strings.SplitN(2) should never return a slice of length > 2
			panic(fmt.Errorf("unexpected only 2 strings to be returned by strings.SplitN but instead got %d", len(parts)))
		}
	}
	return flags, args[i:]
}

func Index[T any](slice []T, i int) (v T, ok bool) {
	if i >= len(slice) {
		return
	}
	return slice[i], true
}
