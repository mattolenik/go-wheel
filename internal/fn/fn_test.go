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
	assert.Equal(2, val)

	ok, val = FindP([]int{1, 2, 3}, func(i *int) bool { return *i == 4 })
	assert.False(ok)
	assert.Equal(0, val)

	ok, val = FindP([]int{}, func(i *int) bool { return *i == 4 })
	assert.False(ok)
	assert.Equal(0, val)
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
