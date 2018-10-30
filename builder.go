package stream

import "reflect"

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

// Repeat
func Repeat(e interface{}, times int) *Stream {
	if times < 1 {
		panic("Stream.Repeat times must greater than 0 ")
	}
	slice := makeReflectSlice(reflect.TypeOf(e), times, times)
	for times > 0 {
		slice.Index(times - 1).Set(reflect.ValueOf(e))
		times--
	}
	return Of(slice.Interface())
}

// RangeInt [from,to)
func RangeInt(from, to int) *Stream {
	if from > to {
		panic("Stream.Repeat from must less than to ")
	}
	slice := make([]int, to-from)
	for i, j := 0, from; i < to; i, j = i+1, j+1 {
		slice[i] = j
	}
	return Of(slice)
}
