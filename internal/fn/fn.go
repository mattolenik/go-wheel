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

func Find[T any](items []T, predicate func(T) bool) (found bool, value T) {
	for _, i := range items {
		if predicate(i) {
			return true, i
		}
	}
	return
}

func FindP[T any](items []T, predicate func(*T) bool) (found bool, value *T) {
	for i := range items {
		if predicate(&items[i]) {
			return true, &items[i]
		}
	}
	return
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

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(item T) bool {
	if _, ok := s[item]; ok {
		return false
	}
	s[item] = struct{}{}
	return true
}

func (s Set[T]) Remove(item T) bool {
	_, ok := s[item]
	delete(s, item)
	return ok
}

func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))
	for item := range s {
		values = append(values, item)
	}
	return values
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
