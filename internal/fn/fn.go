package fn

// TODO: Make some kind of MapE, where function f can return (R, error)?
// Should Map use the ok pattern, returning false when items is false?

func Map[T, R any](items []T, f func(T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = f(item)
	}
	return result
}

func MapP[T, R any](items []T, f func(*T) R) []R {
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = f(&item)
	}
	return result
}

func Find[T any](items []T, predicate func(T) bool) (value T, found bool) {
	for _, i := range items {
		if predicate(i) {
			return i, true
		}
	}
	return
}

func FindP[T any](items []T, predicate func(*T) bool) (value *T, found bool) {
	for i := range items {
		if predicate(&items[i]) {
			return &items[i], true
		}
	}
	return
}

func Has[T any](items []T, predicate func(T) bool) bool {
	_, ok := Find(items, predicate)
	return ok
}

func Filter[T any](items []T, predicate func(*T) bool) []*T {
	result := []*T{}
	for i := range items {
		if predicate(&items[i]) {
			result = append(result, &items[i])
		}
	}
	return result
}

func LookupOrDefault[K comparable, V any](lookupFunc func(key K) (value V, ok bool), key K, defaultVal V) V {
	v, ok := lookupFunc(key)
	if ok {
		return v
	}
	return defaultVal
}

type MultiMap[K comparable, V comparable] map[K]Set[V]

func (m MultiMap[K, V]) Get(key K) Set[V] {
	s, _ := m.Lookup(key)
	return s
}

func (m MultiMap[K, V]) Lookup(key K) (Set[V], bool) {
	v, ok := m[key]
	if ok {
		return v, true
	}
	return Set[V]{}, false
}

func (m MultiMap[K, V]) Put(key K, value V) {
	if _, ok := m[key]; !ok {
		m[key] = Set[V]{}
	}
	m[key].Add(value)
}

func Ptr[T any](v T) *T {
	return &v
}

func Index[T any](slice []T, i int) (v T, ok bool) {
	if i >= len(slice) {
		return
	}
	return slice[i], true
}
