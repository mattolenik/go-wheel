package fn

func Delete[T comparable](items []T, item T) []T {
	i := IndexOf(items, item)
	if i < 0 || items == nil {
		return items
	}
	return append(items[:i], items[i+1:]...)
}

func IndexOf[T comparable](items []T, itemToFind T) int {
	for i, item := range items {
		if item == itemToFind {
			return i
		}
	}
	return -1
}
