package wheel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestSet(t *testing.T) {
	assert := assert.New(t)
	s := Set[int]{}
	assert.True(s.Add(10))
	assert.False(s.Add(10))
	assert.True(s.Add(11))
	assert.ElementsMatch([]int{10, 11}, s.Values())
	assert.True(s.Remove(10))
	assert.ElementsMatch([]int{11}, s.Values())
	assert.False(s.Remove(10))
	assert.ElementsMatch([]int{11}, s.Values())
}
