package pipe

import (
	"fmt"
	"reflect"
)

type Pipe struct {
	elements []interface{}
	src      int
}

func newPipe(length int) *Pipe {
	p := &Pipe{}
	p.elements = make([]interface{}, length)
	return p
}

type Transformer interface {
	Map(fn func(v interface{}) interface{}) (pipe *Pipe)
	Reduce(initialValue interface{}, fn func(prev, v interface{}) interface{}) interface{}
}

// Of create a pipe
func Of(elements interface{}) *Pipe {
	var p *Pipe
	rValue := reflect.ValueOf(elements)
	if rValue.Kind() != reflect.Array || rValue.Kind() != reflect.Slice {
		panic("of func error parameters only support array or slice")
	}
	if rValue.Kind() == reflect.Array || rValue.Kind() == reflect.Slice {
		len := rValue.Len()
		p = newPipe(len)
		for i := 0; i < len; i++ {
			p.elements[i] = rValue.Index(i).Interface()
		}
	}
	return p
}

// Range [start,end)
func Range(start, end int) *Pipe {
	if start < 0 {
		panic(fmt.Sprintf("range func error parameters start[%d] must gather than 0 ", start))
	}
	if end-start >= 0 {
		panic(fmt.Sprintf("range func error parameters start[%d] must less than end[%d] ", start, end))
	}
	var p = newPipe(end - start)
	index := 0
	for i := start; i < end; i++ {
		p.elements[index] = start
		index++
	}
	return p
}

// Map
func (pipe *Pipe) Map(fn func(e interface{}) interface{}) *Pipe {
	for i, v := range pipe.elements {
		pipe.elements[i] = fn(v)
	}
	return pipe
}

// Reduce initialValue 初始值
func (pipe *Pipe) Reduce(initialValue interface{}, fn func(prev, e interface{}) interface{}) interface{} {
	prevValue := initialValue
	for _, v := range pipe.elements {
		prevValue = fn(prevValue, v)
	}
	return prevValue
}
