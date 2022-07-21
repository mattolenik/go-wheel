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
	for _, i := range items {
		if predicate(&i) {
			return true, &i
		}
	}
	return
}

func Ptr[T any](v T) *T {
	return &v
}

func LookupOrDefault[K comparable, V any](lookupFunc func(key K) (value V, ok bool), key K, defaultVal V) V {
	v, ok := lookupFunc(key)
	if ok {
		return v
	}
	return defaultVal
}
