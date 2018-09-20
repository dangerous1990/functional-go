package stream

import "reflect"

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
