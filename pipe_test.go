package pipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	result := Of([]int{1, 2, 3, 4, 5}).Map(func(v interface{}) interface{} {
		return v.(int) * 2
	})
	for i, x := range result.elements {
		assert.Equal(t, (i+1)*2, x)
	}
}
func TestRange(t *testing.T) {
	num := 1
	Range(1, 10).Each(func(i int, x interface{}) {
		assert.Equal(t, x, num)
		num++
	})
}
func TestMap(t *testing.T) {
	index := 1
	Of([]int{1, 2, 3, 4, 5}).Map(func(v interface{}) interface{} {
		assert.Equal(t, index, v)
		index++
		return v
	})
}
func TestMapTo(t *testing.T) {
	Of([]int{1, 2, 3, 4, 5}).MapTo(10).Each(func(i int, e interface{}) {
		assert.Equal(t, 10, e)
	})
}
func TestRepeat(t *testing.T) {
	Repeat(1, 10).Each(func(i int, v interface{}) {
		assert.Equal(t, 1, v)
	})
}
func TestCount(t *testing.T) {
	assert.Equal(t, 10, Repeat(1, 10).Count())
}

func TestFilter(t *testing.T) {
	count := Repeat(1, 10).Filter(func(i int, v interface{}) bool {
		return i%2 == 1
	}).Count()
	assert.Equal(t, 5, count)
}
func TestFirst(t *testing.T) {
	first := Range(1, 10).First()
	assert.Equal(t, 1, first)
}
func TestLast(t *testing.T) {
	last := Range(1, 10).Last()
	assert.Equal(t, 9, last)
}
func TestReverse(t *testing.T) {
	num := 9
	Range(1, 10).Reverse().Each(func(i int, x interface{}) {
		assert.Equal(t, num, x)
		num--
	})
}
