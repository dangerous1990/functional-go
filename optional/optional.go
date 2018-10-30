package optional

import (
	"reflect"
)

type Optional struct {
	value interface{}
}

func Empty() *Optional {
	return &Optional{nil}
}

// Of
func Of(value interface{}) *Optional {
	return &Optional{value}
}

// Get
func (op *Optional) Get() interface{} {
	if op.IsPresent() {
		return op.value
	}
	panic("value is not present")
}

// OrElse
func (op *Optional) OrElse(value interface{}) interface{} {
	if op.IsPresent() {
		return op.value
	}
	return value
}

// OrElseGet
func (op *Optional) OrElseGet(fn func() interface{}) interface{} {
	if op.IsPresent() {
		return op.value
	}
	return reflect.ValueOf(fn).Call([]reflect.Value{})[0].Interface()
}

// IfPresent do something
func (op *Optional) IfPresent(fn func(value interface{})) {
	if op.IsPresent() {
		fnValue := reflect.ValueOf(fn)
		fnValue.Call([]reflect.Value{reflect.ValueOf(op.value)})
	}
}

//  IsPresent
func (op *Optional) IsPresent() bool {
	return !isNil(op.value)
}

func isNil(value interface{}) bool {
	return value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil())
}
