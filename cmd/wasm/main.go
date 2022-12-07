package main

import (
	"fmt"
	"log"
	"strings"
	"syscall/js"
	"time"

	"github.com/jghiloni/aoc2022/pkg/colorize"
	"github.com/jghiloni/aoc2022/pkg/exercise"
	"github.com/jghiloni/aoc2022/pkg/inputs"
	"github.com/jghiloni/aoc2022/pkg/version"
	"github.com/jghiloni/aoc2022/pkg/wasm"
)

func getExercises() js.Value {
	exercises := exercise.ListRegistered()

	infos := make([]any, len(exercises))
	for i := range infos {
		infos[i] = map[string]any{
			"name":  strings.Replace(exercises[i], "day", "Day ", 1),
			"value": exercises[i],
		}
	}

	return js.ValueOf(infos)
}

func runExercise() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		exerciseName := args[0].String()
		e, ok := exercise.GetExercise(exerciseName)
		if !ok {
			return fmt.Errorf("could not find exercise called %q", exerciseName)
		}

		part := args[1].String()

		var worker exercise.ExercisePart
		switch part {
		case "1":
			worker = e.Part1
		case "2":
			worker = e.Part2
		default:
			return fmt.Errorf("invalid exercise part %q", part)
		}

		delay := 100 * time.Millisecond
		if len(args) == 3 {
			dur := args[2].String()

			var err error
			if delay, err = time.ParseDuration(dur); err != nil {
				fmt.Printf("invalid duration %s\n", dur)
			}
		}

		out := colorize.NewColorWriter(wasm.NewWriter(delay), colorize.NewHTMLColorizer())
		output := log.New(out,
			fmt.Sprintf(`{{ colorize "bold;blue" "[%s:part%s] " }}`, exerciseName, part),
			log.Ltime)

		input, err := inputs.Exercises.Open(exerciseName + ".txt")
		if err != nil {
			return fmt.Errorf("error collecting input: %w", err)
		}

		answer, err := worker(input, output)
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}

		return map[string]any{
			"answer": answer,
			"error":  errStr,
		}
	})
}

func waitForever() {
	for {
		<-make(chan bool)
	}
}

func main() {
	fmt.Println("Go WebAssembly")

	js.Global().Set("aocVersion", js.ValueOf(version.Version))
	js.Global().Set("exercises", getExercises())
	js.Global().Set("runExercise", runExercise())

	waitForever()

	log.Println("Exiting")
}
