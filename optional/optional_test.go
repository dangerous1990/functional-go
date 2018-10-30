package optional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	op := Of("test")
	assert.True(t, op.IsPresent())
	assert.Equal(t, "test", op.Get().(string))
}

func TestIfPresent(t *testing.T) {
	op := Of("test")
	op.IfPresent(func(value interface{}) {
		assert.Equal(t, "test", value.(string))
	})
}

func TestOrElse(t *testing.T) {
	op := Empty()
	value := op.OrElse("haha").(string)
	assert.Equal(t, "haha", value)
}

func TestOrElseGet(t *testing.T) {
	op := Empty()
	value := op.OrElseGet(func() interface{} {
		return "test"
	})
	assert.Equal(t, "test", value)
}
