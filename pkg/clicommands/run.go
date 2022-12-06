//go:build unix

package clicommands

import (
	"errors"
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/jghiloni/aoc2022/pkg/colorize"
	"github.com/jghiloni/aoc2022/pkg/exercise"
	"github.com/jghiloni/aoc2022/pkg/inputs"
	"github.com/nexidian/gocliselect"
)

type RunCommand struct {
	Exercise string `short:"e" description:"Which day's exercise to run. Interactive if not specified or invalid"`
	Part     int    `short:"p" description:"Which part of the exercise to run (1 or 2). Interactive if not specified or invalid"`
}

func (r *RunCommand) Run(kCtx *kong.Context) error {
	var (
		e     exercise.Exercise
		part  int
		found bool
	)

	if r.Exercise != "" {
		if e, found = exercise.GetExercise(r.Exercise); !found {
			fmt.Fprintf(kCtx.Stderr, "Could not find exercise %s, launching interactive mode", r.Exercise)
		}
	}

	if e == nil {
		r.Exercise, e = interactiveGetExercise()
	}

	part = r.Part
	if part != 1 && part != 2 {
		part = interactiveGetPart()
	}

	if e == nil {
		return errors.New("could not find exercise")
	}

	var exerciseFunc exercise.ExercisePart
	switch part {
	case 1:
		exerciseFunc = e.Part1
	case 2:
		exerciseFunc = e.Part2
	default:
		return fmt.Errorf("invalid part %d", part)
	}

	inputData, err := inputs.Exercises.Open(fmt.Sprintf("%s.txt", r.Exercise))
	if err != nil {
		return fmt.Errorf("could not open input for %s: %w", r.Exercise, err)
	}
	defer inputData.Close()

	out := colorize.NewColorWriter(kCtx.Stdout, colorize.NewANSIColorizer())
	output := log.New(out,
		fmt.Sprintf(`{{ colorize "bold;blue" "[%s:part%d] " }}`, r.Exercise, part),
		log.Ltime)

	answer, err := exerciseFunc(inputData, output)
	if err != nil {
		return fmt.Errorf("error running exercise %s part %d: %w", r.Exercise, part, err)
	}

	fmt.Fprintln(kCtx.Stdout, "ANSWER: ", answer)
	return nil
}

func interactiveGetExercise() (string, exercise.Exercise) {
	name := selectChoice("Choose an exercise", exercise.ListRegistered()...)
	e, _ := exercise.GetExercise(name)
	return name, e
}

func interactiveGetPart() int {
	partName := selectChoice("Which part?", "Part 1", "Part 2")

	var part int
	fmt.Sscanf(partName, "Part %d", &part)

	return part
}

func selectChoice(prompt string, items ...string) string {
	menu := gocliselect.NewMenu(prompt)
	for _, item := range items {
		menu.AddItem(item, item)
	}

	return menu.Display()
}
