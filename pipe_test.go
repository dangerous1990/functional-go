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
func TestMap(t *testing.T) {
	index := 1
	Of([]int{1, 2, 3, 4, 5}).Map(func(v interface{}) interface{} {
		assert.Equal(t, index, v)
		index++
		return v
	})
}
