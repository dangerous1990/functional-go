package stream

import (
	"reflect"
	"sort"
)

var (
	intType  = reflect.TypeOf(0)
	boolType = reflect.TypeOf(true)
)

const EmptyElement string = "stream_nil"

type transformer interface {
	Map(fn interface{}) *Stream
	Reduce(initialValue interface{}, fn interface{}) interface{}
	Filter(fn interface{}) *Stream
	Sort(lessFunc interface{}) *Stream
	Reverse() *Stream
	First() interface{}
	Last() interface{}
	Each(fn interface{})
}

// Map fn support input one or two args eg. func(i,v int) int   func(v int) int
// fn must be return a value
func (stream *Stream) Map(fn interface{}) *Stream {
	fnValue := reflect.ValueOf(fn)
	fnType := reflect.TypeOf(fn)
	if fnType.NumOut() < 1 {
		panic("Stream.Map(fn), fn is invalid func, must be return a value")
	}
	if !isRightFunc(fnType, []reflect.Type{intType, stream.elementType}, []reflect.Type{fnType.Out(0)}) {
		panic("Stream.Map(fn), fn is invalid func")
	}
	resultSlice := reflect.MakeSlice(reflect.SliceOf(fnType.Out(0)), 0, stream.Length())
	for i := 0; i < stream.Length(); i++ {
		if fnType.NumIn() == 1 {
			resultSlice = reflect.Append(resultSlice, fnValue.Call([]reflect.Value{stream.sourceValue.Index(i)})[0])
		}
		if fnType.NumIn() == 2 {
			resultSlice = reflect.Append(resultSlice, fnValue.Call([]reflect.Value{reflect.ValueOf(i), stream.sourceValue.Index(i)})[0])
		}
	}
	return Of(resultSlice.Interface())
}

// Reduce fn eg. func(prev,v int) int
// prev type must be equal return type
func (stream *Stream) Reduce(initialValue interface{}, fn interface{}) interface{} {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)
	if !isRightFunc(fnType, []reflect.Type{reflect.TypeOf(initialValue), stream.elementType}, []reflect.Type{reflect.TypeOf(initialValue)}) {
		panic("Stream.Reduce(fn), fn is invalid func")
	}
	initValue := reflect.ValueOf(initialValue)
	for i := 0; i < stream.Length(); i++ {
		initValue = fnValue.Call([]reflect.Value{initValue, stream.sourceValue.Index(i)})[0]
	}
	return initValue.Interface()
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

// Sort  lessFunc eg. func(a,b int) bool
func (stream *Stream) Sort(lessFunc interface{}) *Stream {
	lessFuncValue := reflect.ValueOf(lessFunc)
	lessFuncType := lessFuncValue.Type()
	if !isRightFunc(lessFuncType, []reflect.Type{stream.elementType, stream.elementType}, []reflect.Type{boolType}) {
		panic("Stream.Sort(fn) lessFunc is invalid func")
	}
	delegate := &sortDelegate{
		Arr:      reflect.ValueOf(stream.Get()),
		lessFunc: lessFuncValue,
	}
	// sort
	sort.Stable(delegate)
	return Of(delegate.Arr.Interface())
}

// Reverse
func (stream *Stream) Reverse() *Stream {
	slice := reflect.MakeSlice(stream.sliceType, stream.Length(), stream.Length())
	for i, j := 0, stream.Length()-1; i <= j; i, j = i+1, j-1 {
		slice.Index(i).Set(stream.sourceValue.Index(j))
		slice.Index(j).Set(stream.sourceValue.Index(i))
	}
	return Of(slice.Interface())
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

// Each
func (stream *Stream) Each(fn interface{}) {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)
	if !isRightFunc(fnType, []reflect.Type{intType, stream.elementType}, []reflect.Type{}) {
		panic("Stream.Filter(fn), fn is invalid func")
	}
	for i := 0; i < stream.Length(); i++ {
		elementValue := stream.sourceValue.Index(i)
		if fnType.NumIn() == 1 {
			fnValue.Call([]reflect.Value{elementValue})
		}
		if fnType.NumIn() == 2 {
			fnValue.Call([]reflect.Value{reflect.ValueOf(i), elementValue})
		}
	}
}

func isTypeMatched(a, b reflect.Type) bool {
	if a == b {
		return true
	}
	if a.Kind() == reflect.Interface {
		return b.Implements(a)
	}
	return false
}

func isRightFunc(funcType reflect.Type, inputTypes, outputTypes []reflect.Type) bool {
	if funcType.Kind() != reflect.Func {
		return false
	}
	if funcType.NumIn() > len(inputTypes) {
		return false
	}
	if funcType.NumOut() != len(outputTypes) {
		return false
	}
	if funcType.NumOut() > 0 && !isTypeMatched(funcType.Out(0), outputTypes[0]) {
		return false
	}
	if funcType.NumIn() == len(inputTypes) {
		for i, t := range inputTypes {
			if !isTypeMatched(t, funcType.In(i)) {
				return false
			}
		}
	}
	if funcType.NumIn() < len(inputTypes) {
		if !isTypeMatched(funcType.In(0), inputTypes[len(inputTypes)-1]) {
			return false
		}
	}

	return true
}
