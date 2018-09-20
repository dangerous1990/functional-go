package stream

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

func TestRepeat(t *testing.T) {
	result := Repeat(2, 10).Get().([]int)
	for _, v := range result {
		assert.Equal(t, 2, v)
	}
}

func TestRepeatWithStruct(t *testing.T) {
	type square struct {
		x int
		y int
	}
	result := Repeat(&square{
		10,
		20,
	}, 10).Get().([]*square)
	for _, v := range result {
		assert.Equal(t, 10, v.x)
		assert.Equal(t, 20, v.y)
	}
}

func TestRangeInt(t *testing.T) {
	RangeInt(0, 10).Each(func(i, v int) {
		assert.Equal(t, i, v)
	})

}
