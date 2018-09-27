package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	even := Of([]int{1, 2, 3, 4, 5}).Filter(func(i, x int) bool {
		return x%2 == 0
	}).Get().([]int)
	odd := Of([]int{1, 2, 3, 4, 5}).Filter(func(x int) bool {
		return x%2 == 1
	}).Get().([]int)

	for _, v := range even {
		assert.Equal(t, v%2, 0)
	}

	for _, v := range odd {
		assert.Equal(t, v%2, 1)
	}
}

func TestFirstLast(t *testing.T) {
	first := Of([]int{0, 1, 2, 3, 4, 5}).First().(int)
	last := Of([]int{0, 1, 2, 3, 4, 5}).Last().(int)
	assert.Equal(t, 0, first)
	assert.Equal(t, 5, last)
}

func TestSkip(t *testing.T) {
	x := 3
	Of([]int{0, 1, 2, 3, 4, 5}).Skip(3).Each(func(i, v int) {
		assert.Equal(t, x, v)
		x++
	})
}

func TestFind(t *testing.T) {
	result := Of([]int{1, 2, 3, 4, 5}).Find(func(v int) bool {
		return v%2 == 0
	}).(int)
	assert.Equal(t, 2, result)
}
func TestFindIndex(t *testing.T) {
	result := Of([]int{1, 2, 3, 4, 5}).FindIndex(func(v int) bool {
		return v%2 == 0
	})
	assert.Equal(t, 1, result)
}
