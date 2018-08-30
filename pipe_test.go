package pipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	result := Of([]int{1, 2, 3, 4, 5}).Map(func(i, value int) int {
		return value * 2
	}).Get().([]int)
	for i, v := range result {
		assert.Equal(t, (i+1)*2, v)
	}
}
func TestMap(t *testing.T) {
	result := Of([]int{1, 2, 3, 4, 5}).Map(func(i, value int) int {
		return value * 2
	}).Map(func(i, value int) int {
		return value / 2
	}).Get().([]int)
	for i, v := range result {
		assert.Equal(t, i+1, v)
	}
}

func TestReduce(t *testing.T) {
	sum := Of([]int{1, 2, 3, 4, 5}).Reduce(0, func(prev, v int) int {
		return prev + v
	}).(int)
	assert.Equal(t, 15, sum)
}

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
func TestSort(t *testing.T) {
	result := Of([]int{0, 2, 5, 1, 4, 3}).Sort(func(a, b int) bool {
		return a < b
	}).Get().([]int)
	for i, v := range result {
		assert.Equal(t, i, v)
	}
}
