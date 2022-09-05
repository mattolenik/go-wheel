package twine

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mattolenik/go-wheel/pkg/typ"
)

// TODO: Add tests for this package

// TODO: Better name
// TODO: Move to another package?
func FromDelimetedList[T typ.StringRepresentable](str, delim string) ([]T, error) {
	parts := strings.Split(str, delim)
	if len(parts) == 0 || (len(parts) == 1 || parts[0] == "") {
		return []T{}, nil
	}
	result := make([]T, len(parts))
	resultValue := reflect.ValueOf(result)
	for i, part := range parts {
		part = strings.TrimSpace(part)
		v, err := typ.Parse[T](part)
		if err != nil {
			return nil, fmt.Errorf("list item %q is not a valid value: %w", part, err)
		}
		resultValue.Index(i).Set(reflect.ValueOf(v))
	}
	return result, nil
}
