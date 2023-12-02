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
	TypedOption any
	Global      bool
	Default     *string
}

type TypedOption[T CommandLineType] struct {
	Option
	Value *T
}

func (o *TypedOption[T]) Bind(value *T) *TypedOption[T] {
	o.Value = value
	return o
}

func (o *TypedOption[T]) WithDefault(val string) *TypedOption[T] {
	o.Default = &val
	return o
}

func (o *TypedOption[T]) Required() *TypedOption[T] {
	o.IsRequired = true
	return o
}
