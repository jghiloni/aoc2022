package utils

type Stack[T any] struct {
	s []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		s: []T{},
	}
}

func (s *Stack[T]) Push(item T) int {
	s.s = append([]T{item}, s.s...)
	return len(s.s)
}

func (s *Stack[T]) Pop() (T, int) {
	if len(s.s) == 0 {
		tp := new(T)
		return *tp, -1
	}

	r := s.s[0]
	switch len(s.s) {
	case 1:
		s.s = []T{}
	default:
		s.s = s.s[1:]
	}

	return r, len(s.s)
}

func (s *Stack[T]) Peek() T {
	if len(s.s) == 0 {
		tp := new(T)
		return *tp
	}

	return s.s[0]
}

func (s *Stack[T]) Size() int {
	return len(s.s)
}

func (s *Stack[T]) ToFIFOSlice() []T {
	retVal := make([]T, len(s.s))
	for i := range s.s {
		retVal[len(s.s)-i-1] = s.s[i]
	}

	return retVal
}
