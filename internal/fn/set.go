package fn

// TODO: Should Set and other mutable data structures be moved into a separate package for mutable types?
//     : Should there be a separate package for slice-like "immutable" types that are more like native Go collections?

type Set[T comparable] map[T]struct{}

// Add adds an item to the set. If the item does not exist and gets added, Add will return true. If it's already
// present, no operation will occur and Add will return false.
func (s Set[T]) Add(item T) bool {
	if _, has := s[item]; has {
		return false
	}
	s[item] = struct{}{}
	return true
}

func (s Set[T]) Remove(item T) bool {
	_, has := s[item]
	delete(s, item)
	return has
}

func (s Set[T]) Values() []T {
	values := make([]T, 0, len(s))
	for item := range s {
		values = append(values, item)
	}
	return values
}

type OrderedSet[T comparable] struct {
	set     Set[T]
	ordered []T
}

// lazy initializes defaults for this struct.
// This awkward pattern allows the user to not need to use
// TODO: really? below
// a constructor method, they can just write: s := OrderedSet[T]{1,2,3}
func (s *OrderedSet[T]) lazy() {
	if s.set == nil {
		s.set = Set[T]{}
	}
	if s.ordered == nil {
		s.ordered = []T{}
	}
}

func (s *OrderedSet[T]) Add(item T) bool {
	s.lazy()
	added := s.set.Add(item)
	if !added {
		return false
	}
	s.ordered = append(s.ordered, item)
	return true
}

func (s *OrderedSet[T]) AddAll(items ...T) bool {
	s.lazy()
	var added bool
	for _, item := range items {
		a := s.Add(item)
		added = added || a
	}
	return added
}

func (s *OrderedSet[T]) Remove(item T) bool {
	s.lazy()
	removed := s.set.Remove(item)
	if !removed {
		return false
	}
	s.ordered = Delete(s.ordered, item)
	return true
}

func (s *OrderedSet[T]) Values() []T {
	s.lazy()
	return s.ordered
}

func (s *OrderedSet[T]) ValuesP() []*T {
	result := make([]*T, len(s.ordered))
	for i, item := range s.ordered {
		result[i] = &item
	}
	return result
}

func (s *OrderedSet[T]) Contains(item T) bool {
	_, has := s.set[item]
	return has
}

func (s *OrderedSet[T]) ContainsP(item *T) bool {
	_, has := s.set[*item]
	return has
}

func (s *OrderedSet[T]) ContainsAll(items []T) bool {
	for _, item := range items {
		if !s.Contains(item) {
			return false
		}
	}
	return true
}

func (s *OrderedSet[T]) ContainsAny(items []T) bool {
	for _, item := range items {
		if s.Contains(item) {
			return true
		}
	}
	return false
}

func (s *OrderedSet[T]) ContainsAnyP(items []*T) bool {
	for _, item := range items {
		if s.ContainsP(item) {
			return true
		}
	}
	return false
}

func (s *OrderedSet[T]) ContainsAllP(items []*T) bool {
	for _, item := range items {
		if !s.ContainsP(item) {
			return false
		}
	}
	return true
}
