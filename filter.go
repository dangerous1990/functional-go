package stream

import (
	"reflect"
)

type filter interface {
	Filter(fn interface{}) *Stream
	First() interface{}
	Last() interface{}
	Skip(n int) interface{}
	Find(fn interface{}) interface{}
}

// Filter fn support input one or two args eg. func(i,v int) int   func(v int) int
// fn must be return bool
func (stream *Stream) Filter(fn interface{}) *Stream {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)
	if !isRightFunc(fnType, []reflect.Type{intType, stream.elementType}, []reflect.Type{boolType}) {
		panic("Stream.Filter(fn), fn is invalid func")
	}
	resultSlice := reflect.MakeSlice(stream.sliceType, 0, stream.Length())
	for i := 0; i < stream.Length(); i++ {
		elementValue := stream.sourceValue.Index(i)
		if fnType.NumIn() == 1 {
			if fnValue.Call([]reflect.Value{elementValue})[0].Bool() {
				resultSlice = reflect.Append(resultSlice, elementValue)
			}
		}
		if fnType.NumIn() == 2 {
			if fnValue.Call([]reflect.Value{reflect.ValueOf(i), elementValue})[0].Bool() {
				resultSlice = reflect.Append(resultSlice, elementValue)
			}
		}
	}
	return Of(resultSlice.Interface())
}

// First
func (stream *Stream) First() interface{} {
	if stream.Length() < 1 {
		return EmptyElement
	}
	return stream.sourceValue.Index(0).Interface()
}

// Last
func (stream *Stream) Last() interface{} {
	if stream.Length() < 1 {
		return EmptyElement
	}
	return stream.sourceValue.Index(stream.Length() - 1).Interface()
}

// Skip
func (stream *Stream) Skip(n int) *Stream {
	if stream.Length() < n {
		panic("stream.Skip n must less than stream.Length()")
	}
	slice := reflect.MakeSlice(stream.sliceType, stream.Length()-n, stream.Length()-n)
	for i, j := 0, n; j < stream.Length(); i, j = i+1, j+1 {
		slice.Index(i).Set(stream.sourceValue.Index(j))
	}
	return Of(slice.Interface())
}

// Find
func (stream *Stream) Find(fn interface{}) interface{} {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)
	if !isRightFunc(fnType, []reflect.Type{intType, stream.elementType}, []reflect.Type{boolType}) {
		panic("Stream.Find(fn), fn is invalid func")
	}
	for i := 0; i < stream.Length(); i++ {
		elementValue := stream.sourceValue.Index(i)
		if fnType.NumIn() == 1 {
			if fnValue.Call([]reflect.Value{elementValue})[0].Bool() {
				return elementValue.Interface()
			}
		}
		if fnType.NumIn() == 2 {
			if fnValue.Call([]reflect.Value{reflect.ValueOf(i), elementValue})[0].Bool() {
				return elementValue.Interface()
			}
		}
	}
	return nil
}
