package wheel

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
