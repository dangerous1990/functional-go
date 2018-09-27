package stream

import (
	"fmt"
	"reflect"
)

type optional struct {
	value interface{}
}

var Optional *optional = defaultOptional()

func defaultOptional() *optional {
	return &optional{}
}

func (op *optional) Of(value interface{}) *optional {
	return &optional{value}
}

func (op *optional) IsPresent() bool {
	return !isNil(op.value)
}

// TODO update
func isNil(value interface{}) bool {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("panic in isNil(), err:%v\n", e)
		}
	}()
	if value == nil {
		return true
	}
	if reflect.ValueOf(value).IsNil() {
		return true
	}
	return false
}
