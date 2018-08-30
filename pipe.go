package pipe

import (
	"reflect"
)

var INT_TYPE reflect.Type = reflect.TypeOf(int(0))

type Stream struct {
	parentStream *Stream
	source       interface{}  // source data it must be array or slice
	sliceType    reflect.Type // source data sliceType
	elementType  reflect.Type // source data element type
	sourceType   reflect.Type
	sourceValue  reflect.Value
}
type FuncDesp struct {
	inputType  []reflect.Type
	outputType reflect.Type // only one return field
	funcType   reflect.Type
}
type Transformer interface {
	Map(fn interface{}) *Stream
	Reduce(initialValue interface{}, fn interface{}) interface{}
	// MapTo(v interface{}) *Stream
	// Reverse() *Stream
	// Each(fn func(i int, v interface{}))
	// Filter(fn func(i int, v interface{})) *Stream
	// First() interface{}
	// Last() interface{}
	// IsEmpty() bool
	// Count() int
	// Max()
	// Get() interface{}
}

func (stream *Stream) Length() int {
	return reflect.ValueOf(stream.source).Len()
}
func newStream(source interface{}) *Stream {
	sourceValue := reflect.ValueOf(source)
	sourceType := reflect.TypeOf(source)
	if !(sourceValue.Kind() == reflect.Array || sourceValue.Kind() == reflect.Slice) {
		panic("of func error parameters only support array or slice")
	}
	elementType := sourceType.Elem()
	sliceType := reflect.SliceOf(elementType)
	return &Stream{
		parentStream: nil,
		source:       source,
		sliceType:    sliceType,
		elementType:  elementType,
		sourceType:   sourceType,
		sourceValue:  sourceValue,
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
	// if !isTypeMatched(funcType.Out(0), outputTypes[0]) {
	// 	return false
	// }
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

// Of create a Stream
func Of(source interface{}) *Stream {
	return newStream(source)
}

// Map
func (stream *Stream) Map(fn interface{}) *Stream {
	fnValue := reflect.ValueOf(fn)
	fnType := reflect.TypeOf(fn)
	if !isRightFunc(fnType, []reflect.Type{INT_TYPE, stream.elementType}, []reflect.Type{nil}) {
		panic("map invalid func")
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
	return newStream(resultSlice.Interface())
}

// ToSlice
func (stream *Stream) ToSlice() interface{} {
	return stream.source
}

// Reduce initialValue 初始值
func (stream *Stream) Reduce(initialValue interface{}, fn interface{}) interface{} {
	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)
	if !isRightFunc(fnType, []reflect.Type{reflect.TypeOf(initialValue), stream.elementType}, []reflect.Type{reflect.TypeOf(initialValue)}) {
		panic("wrong reduce func")
	}
	initValue := reflect.ValueOf(initialValue)
	for i := 0; i < stream.Length(); i++ {
		initValue = fnValue.Call([]reflect.Value{initValue, stream.sourceValue.Index(i)})[0]
	}
	return initValue.Interface()
}
