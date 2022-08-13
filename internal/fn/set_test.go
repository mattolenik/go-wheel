package fn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestOrderedSet(t *testing.T) {
	assert := assert.New(t)

	{
		s := OrderedSet[int]{}
		s.AddAll()
		assert.Equal([]int{}, s.Values())
		assert.True(s.Add(1))
		assert.True(s.Remove(1))

		assert.Equal([]int{}, s.Values())
		assert.True(s.Add(1))
		assert.False(s.Remove(2))
		assert.Equal([]int{1}, s.Values())
	}
	{
		s := OrderedSet[int]{}
		assert.True(s.AddAll(1, 2, 3))
		assert.Equal([]int{1, 2, 3}, s.Values())
	}
	{
		s := OrderedSet[int]{}
		assert.True(s.AddAll(3, 2, 1))
		assert.Equal([]int{3, 2, 1}, s.Values())
	}
	{
		s := OrderedSet[int]{}
		assert.True(s.AddAll(1, 2, 3))
		assert.True(s.Remove(2))
		assert.Equal([]int{1, 3}, s.Values())
	}
	{
		s := OrderedSet[int]{}
		assert.True(s.AddAll(3, 2, 1))
		assert.True(s.Remove(2))
		assert.Equal([]int{3, 1}, s.Values())
	}
	{
		s := OrderedSet[int]{}
		assert.True(s.AddAll(3, 2, 1))
		assert.True(s.Remove(2))
		assert.True(s.Remove(1))
		assert.Equal([]int{3}, s.Values())
	}
}
