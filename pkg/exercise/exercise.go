package exercise

import (
	"io"
	"log"
	"sort"
	"sync"
)

type ExercisePart func(io.Reader, *log.Logger) (any, error)

type Exercise interface {
	Part1(input io.Reader, output *log.Logger) (any, error)
	Part2(input io.Reader, output *log.Logger) (any, error)
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
