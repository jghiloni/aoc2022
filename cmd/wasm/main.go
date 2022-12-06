package main

import (
	"encoding/json"
	"errors"
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

type exerciseInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func getExercises() js.Func {
	exercises := exercise.ListRegistered()

	infos := make([]exerciseInfo, len(exercises))
	for i := range infos {
		infos[i] = exerciseInfo{
			Name:  strings.Replace(exercises[i], "day", "Day ", 1),
			Value: exercises[i],
		}
	}

	return js.FuncOf(func(this js.Value, args []js.Value) any {
		jsonBytes, err := json.MarshalIndent(infos, "", "  ")
		if err != nil {
			fmt.Printf("unable to marshal json: %v\n", err)
			return err.Error()
		}

		return string(jsonBytes)
	})
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

		if !args[2].Truthy() {
			return errors.New("the third argument must be an HTMLElement")
		}

		out := colorize.NewColorWriter(wasm.NewWriter(100*time.Millisecond), colorize.NewHTMLColorizer())
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
	js.Global().Set("getExercises", getExercises())
	js.Global().Set("runExercise", runExercise())

	waitForever()

	log.Println("Exiting")
}
