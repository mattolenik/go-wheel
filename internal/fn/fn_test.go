package fn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	assert := assert.New(t)

	ok, val := Find([]int{1, 2, 3}, func(i int) bool { return i == 2 })
	assert.True(ok)
	assert.Equal(2, val)

	ok, val = Find([]int{1, 2, 3}, func(i int) bool { return i == 4 })
	assert.False(ok)
	assert.Equal(0, val)

	ok, val = Find([]int{}, func(i int) bool { return i == 4 })
	assert.False(ok)
	assert.Equal(0, val)
}

func TestFindP(t *testing.T) {
	assert := assert.New(t)

	ok, val := FindP([]int{1, 2, 3}, func(i *int) bool { return *i == 2 })
	assert.True(ok)
	assert.Equal(2, *val)

	ok, _ = FindP([]int{1, 2, 3}, func(i *int) bool { return *i == 4 })
	assert.False(ok)

	ok, _ = FindP([]int{}, func(i *int) bool { return *i == 4 })
	assert.False(ok)
}

func TestMap(t *testing.T) {
	assert := assert.New(t)
	assert.Equal([]string{"1", "2"}, Map([]int{1, 2}, func(i int) string { return fmt.Sprintf("%d", i) }))
	assert.Equal([]string{}, Map([]int{}, func(i int) string { return "" }))
	assert.Equal([]string{}, Map(nil, func(i int) string { return "" }))
}

func TestMapP(t *testing.T) {
	assert := assert.New(t)
	assert.Equal([]string{"1", "2"}, MapP([]int{1, 2}, func(i *int) string { return fmt.Sprintf("%d", *i) }))
	assert.Equal([]string{}, MapP([]int{}, func(i *int) string { return "" }))
	assert.Equal([]string{}, MapP(nil, func(i *int) string { return "" }))
}

func TestMultiMap(t *testing.T) {
	assert := assert.New(t)

	m := MultiMap[int, int]{}
	m.Put(1, 10)
	m.Put(1, 11)
	assert.ElementsMatch([]int{10, 11}, m.Get(1).Values())
	assert.ElementsMatch([]int{}, m.Get(100).Values())
	v, ok := m.Lookup(1)
	assert.True(ok)
	assert.ElementsMatch([]int{10, 11}, v.Values())
	v, ok = m.Lookup(100)
	assert.False(ok)
	assert.ElementsMatch([]int{}, v.Values())
}

func TestFilter(t *testing.T) {
	assert := assert.New(t)
	items := []int{1, 2, 3, 4}
	expected := []int{3, 4}

	filtered := Filter(items, func(i *int) bool { return *i > 2 })
	assert.ElementsMatch(expected, Map(filtered, func(i *int) int { return *i }))

	assert.ElementsMatch([]*int{}, Filter([]int{}, func(i *int) bool { return true }))
}
