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
	Map(fn func(v interface{}) interface{}) *Pipe
	Reduce(initialValue interface{}, fn func(prev, v interface{}) interface{}) interface{}
	Count() int
	Each(fn func(i int, v interface{}))
	Filter(fn func(i int, v interface{})) *Pipe
	First() interface{}
	Last() interface{}
	IsEmpty() bool
}

// Of create a pipe
func Of(elements interface{}) *Pipe {
	var p *Pipe
	rValue := reflect.ValueOf(elements)
	if !(rValue.Kind() == reflect.Array || rValue.Kind() == reflect.Slice) {
		panic("of func error parameters only support array or slice")
	}
	len := rValue.Len()
	p = newPipe(len)
	for i := 0; i < len; i++ {
		p.elements[i] = rValue.Index(i).Interface()
	}
	return p
}

// Range [start,end)
func Range(start, end int) *Pipe {
	if start < 0 {
		panic(fmt.Sprintf("range func error parameters start[%d] must gather than 0 ", start))
	}
	if end-start < 0 {
		panic(fmt.Sprintf("range func error parameters start[%d] must less than end[%d] ", start, end))
	}
	p := newPipe(end - start)
	index := 0
	for i := start; i < end; i++ {
		p.elements[index] = i
		index++
	}
	return p
}

// Repeat
func Repeat(e interface{}, times int) *Pipe {
	p := newPipe(times)
	for i := 0; i < times; i++ {
		p.elements[i] = e
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

// Count
func (pipe *Pipe) Count() int {
	if pipe.elements == nil {
		return 0
	}
	return len(pipe.elements)
}

// Each
func (pipe *Pipe) Each(fn func(i int, v interface{})) {
	for i, v := range pipe.elements {
		fn(i, v)
	}
}

// Filter
func (pipe *Pipe) Filter(fn func(i int, v interface{}) bool) *Pipe {
	filterElements := make([]interface{}, 0)
	for i, v := range pipe.elements {
		if fn(i, v) {
			filterElements = append(filterElements, v)
		}
	}
	pipe.elements = filterElements
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
func (pipe *Pipe) First() interface{} {
	if pipe.IsEmpty() {
		return nil
	}
	return pipe.elements[0]
}

func (pipe *Pipe) IsEmpty() bool {
	if pipe.elements == nil {
		return true
	}
	return len(pipe.elements) == 0
}
func (pipe *Pipe) Last() interface{} {
	if pipe.IsEmpty() {
		return nil
	}
	return pipe.elements[len(pipe.elements)-1]
}
