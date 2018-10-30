package functional

import (
	"reflect"
)

type Stream struct {
	// source data it must be array or slice
	source interface{}
	// source data sliceType
	sliceType reflect.Type
	// source data element type
	elementType reflect.Type
	sourceType  reflect.Type
	sourceValue reflect.Value
}

// Get get source
func (stream *Stream) Get() interface{} {
	return stream.source
}

// Length get source length
func (stream *Stream) Length() int {
	return stream.sourceValue.Len()
}

// IsEmpty Stream slice is empty
func (stream *Stream) IsEmpty() bool {
	return stream.Length() == 0
}

// GetByIndex get element by index
func (stream *Stream) GetElement(index int) interface{} {
	return stream.sourceValue.Index(index).Interface()
}

// Index
func (stream *Stream) Index(index int) interface{} {
	return stream.getElement(index).Interface()
}

func (stream *Stream) getElement(index int) reflect.Value {
	return stream.sourceValue.Index(index)
}
