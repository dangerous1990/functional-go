package stream

import "reflect"

type builder interface {
	Of(slice interface{}) *Stream
}

func newStream(source interface{}) *Stream {
	if source == nil {
		panic("new stream failed, source is nil")
	}
	sourceValue := reflect.ValueOf(source)
	sourceType := reflect.TypeOf(source)
	if !(sourceValue.Kind() == reflect.Array || sourceValue.Kind() == reflect.Slice) {
		panic("new stream failed, error parameters only support array or slice")
	}
	elementType := sourceType.Elem()
	sliceType := reflect.SliceOf(elementType)
	return &Stream{
		source:      source,
		sliceType:   sliceType,
		elementType: elementType,
		sourceType:  sourceType,
		sourceValue: sourceValue,
	}
}

// Of create a Stream
func Of(source interface{}) *Stream {
	return newStream(source)
}
