package fn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	assert.Equal([]int(nil), Delete(nil, 1))
	assert.Equal([]int{}, Delete([]int{}, 1))
	assert.Equal([]int{}, Delete([]int{1}, 1))
	assert.Equal([]int{1}, Delete([]int{1, 2}, 2))
	assert.Equal([]int{2}, Delete([]int{1, 2}, 1))
	assert.Equal([]int{1, 2}, Delete([]int{1, 2, 3}, 3))
	assert.Equal([]int{1, 3}, Delete([]int{1, 2, 3}, 2))
	assert.Equal([]int{2, 3}, Delete([]int{2, 3}, 1))
}
