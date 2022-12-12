package utils

import (
	"fmt"
	"strings"
)

type Queue[T any] struct {
	s []T
}

func NewQueue[T any](items ...T) *Queue[T] {
	return &Queue[T]{
		s: items,
	}
}

func (s *Queue[T]) Push(item T) int {
	s.s = append(s.s, item)
	return len(s.s)
}

func (s *Queue[T]) Pop() (T, int) {
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

func (s *Queue[T]) Peek() T {
	if len(s.s) == 0 {
		tp := new(T)
		return *tp
	}

	return s.s[0]
}

func (s *Queue[T]) Size() int {
	return len(s.s)
}

func (s *Queue[T]) Join(sep string) string {
	str := ""
	for _, i := range s.s {
		str = fmt.Sprintf("%s%s%v", str, sep, i)
	}

	return strings.TrimPrefix(str, sep)
}
