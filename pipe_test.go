package pipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	result := Of([]int{1, 2, 3, 4, 5}).Map(func(i, value int) int {
		return value * 2
	}).ToSlice().([]int)
	for i, v := range result {
		assert.Equal(t, (i+1)*2, v)
	}
}
func TestMap(t *testing.T) {
	result := Of([]int{1, 2, 3, 4, 5}).Map(func(i, value int) int {
		return value * 2
	}).Map(func(i, value int) int {
		return value / 2
	}).ToSlice().([]int)
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
