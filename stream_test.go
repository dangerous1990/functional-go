package stream

import (
	"math/rand"
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

const num = 1000

func TestGeneralCost(t *testing.T) {
	slice := make([]int, num)
	count := num
	for count > 0 {
		slice = append(slice, rand.Int())
		count--
	}
	for i, v := range slice {
		slice[i] = v * 2
	}

}

func TestFunctionalCost(t *testing.T) {
	slice := make([]int, num)
	count := num
	for count > 0 {
		slice = append(slice, rand.Int())
		count--
	}
	var _ = Of([]int{1, 2, 3, 4, 5}).Map(func(i, value int) int {
		return value * 2
	})

}

func TestReverse(t *testing.T) {
	result := Of([]int{0, 1, 2, 3, 4, 5}).Reverse().Get().([]int)
	index := 5
	for _, v := range result {
		assert.Equal(t, index, v)
		index--
	}
}
func TestEach(t *testing.T) {
	index := 0
	Of([]int{0, 1, 2, 3, 4, 5}).Each(func(v int) {
		assert.Equal(t, index, v)
		index++
	})
}
func TestFirstLast(t *testing.T) {
	first := Of([]int{0, 1, 2, 3, 4, 5}).First().(int)
	last := Of([]int{0, 1, 2, 3, 4, 5}).Last().(int)
	assert.Equal(t, 0, first)
	assert.Equal(t, 5, last)
}
