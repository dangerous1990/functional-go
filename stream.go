package stream

import (
	"reflect"
)

type Stream struct {
	source      interface{}  // source data it must be array or slice
	sliceType   reflect.Type // source data sliceType
	elementType reflect.Type // source data element type
	sourceType  reflect.Type
	sourceValue reflect.Value
	err         error // default is nil ,when operation panic is the recover error
}

// Get get source
func (stream *Stream) Get() interface{} {
	return stream.source
}
