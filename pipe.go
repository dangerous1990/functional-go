package pipe

import (
	"reflect"
	"sort"
)

var (
	intType  = reflect.TypeOf(0)
	boolType = reflect.TypeOf(true)
)

type Stream struct {
	source      interface{}  // source data it must be array or slice
	sliceType   reflect.Type // source data sliceType
	elementType reflect.Type // source data element type
	sourceType  reflect.Type
	sourceValue reflect.Value
}

type transformer interface {
	Map(fn interface{}) *Stream
	Reduce(initialValue interface{}, fn interface{}) interface{}
	Filter(fn interface{}) *Stream
	Sort(lessFunc interface{}) *Stream
	// MapTo(v interface{}) *Stream
	// Reverse() *Stream
	// Each(fn func(i int, v interface{}))
	// First() interface{}
	// Last() interface{}
	// IsEmpty() bool
	// Count() int
	// Max()
	// Get() interface{}
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
	if !isTypeMatched(funcType.Out(0), outputTypes[0]) {
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

type sortDelegate struct {
	Arr      reflect.Value
	lessFunc reflect.Value
}

func (s *sortDelegate) Len() int {
	return s.Arr.Len()
}

func (s *sortDelegate) Less(i, j int) bool {
	outs := s.lessFunc.Call([]reflect.Value{s.Arr.Index(i), s.Arr.Index(j)})
	return outs[0].Interface().(bool)
}

func (s *sortDelegate) Swap(i, j int) {
	ti := s.Arr.Index(i).Interface()
	tj := s.Arr.Index(j).Interface()
	s.Arr.Index(i).Set(reflect.ValueOf(tj))
	s.Arr.Index(j).Set(reflect.ValueOf(ti))
}

// Length get source length
func (stream *Stream) Length() int {
	return stream.sourceValue.Len()
}

// Get get source
func (stream *Stream) Get() interface{} {
	return stream.source
}

// Of create a Stream
func Of(source interface{}) *Stream {
	return newStream(source)
}

// Map fn support input one or two paramter eg. func(i,v int) int   func(v int) int
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
	return newStream(resultSlice.Interface())
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

// Filter fn support input one or two paramter eg. func(i,v int) int   func(v int) int
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
	return newStream(resultSlice.Interface())
}

// Sort  lessFunc eg. func(a,b int) bool
func (stream *Stream) Sort(lessFunc interface{}) *Stream {
	lessFuncValue := reflect.ValueOf(lessFunc)
	lessFuncType := lessFuncValue.Type()
	if !isRightFunc(lessFuncType, []reflect.Type{stream.elementType, stream.elementType}, []reflect.Type{boolType}) {
		panic("sort less function invalid")
	}
	delegate := &sortDelegate{
		Arr:      reflect.ValueOf(stream.Get()),
		lessFunc: lessFuncValue,
	}
	// sort
	sort.Stable(delegate)
	return newStream(delegate.Arr.Interface())
}
