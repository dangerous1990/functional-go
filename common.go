package functional

import (
	"reflect"
)

func makeReflectSlice(typ reflect.Type, len, cap int) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(typ), len, cap)
}

func makeEmptyReflectSlice(typ reflect.Type) reflect.Value {
	return reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
}
