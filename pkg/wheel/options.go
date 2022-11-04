package wheel

import "reflect"

type Option struct {
	Name        string
	Description string
	Examples    []string
	IsRequired  bool
	Type        reflect.Type
	Setter      func(string) error
	Get         func() any
}

type TypedOption[T CommandLineType] struct {
	Option
	Value   *T
	Default *T
}

func (o *TypedOption[T]) Bind(value *T) *TypedOption[T] {
	o.Value = value
	return o
}

func (o *TypedOption[T]) WithDefault(dflt T) *TypedOption[T] {
	o.Default = &dflt
	return o
}

func (o *TypedOption[T]) Required() *TypedOption[T] {
	o.IsRequired = true
	return o
}
