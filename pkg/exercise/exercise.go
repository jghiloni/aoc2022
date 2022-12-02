package exercise

import (
	"io"
	"sort"
	"sync"
)

type ExercisePart func(io.Reader, io.Writer, io.Writer) (any, error)

type Exercise interface {
	Part1(stdin io.Reader, stdout io.Writer, stderr io.Writer) (any, error)
	Part2(stdin io.Reader, stdout io.Writer, stderr io.Writer) (any, error)
}

var exerciseRegistry map[string]Exercise = map[string]Exercise{}
var registryMutex sync.RWMutex = sync.RWMutex{}

func ListRegistered() []string {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	names := []string{}
	for name := range exerciseRegistry {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

func Register(name string, exercise Exercise) {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	exerciseRegistry[name] = exercise
}

func GetExercise(name string) (Exercise, bool) {
	registryMutex.RLock()
	defer registryMutex.RUnlock()

	exercise, ok := exerciseRegistry[name]
	return exercise, ok
}
