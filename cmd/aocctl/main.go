package main

import (
	"github.com/alecthomas/kong"
	"github.com/jghiloni/aoc2022/pkg/clicommands"
)

var cliOptions struct {
	Run      *clicommands.RunCommand              `cmd:"" description:"Run an exercise"`
	Generate *clicommands.GenerateExerciseCommand `cmd:"" description:"Generate a day's code"`
}

func main() {
	kCtx := kong.Parse(&cliOptions)
	kCtx.FatalIfErrorf(kCtx.Run())
}
