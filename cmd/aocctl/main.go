//go:build unix

package main

import (
	"github.com/alecthomas/kong"
	"github.com/jghiloni/aoc2022/pkg/clicommands"
	"github.com/jghiloni/aoc2022/pkg/version"
)

var cliOptions struct {
	Run         *clicommands.RunCommand              `cmd:"" description:"Run an exercise"`
	Generate    *clicommands.GenerateExerciseCommand `cmd:"" description:"Generate a day's code"`
	LocalServer *clicommands.LocalServerCommand      `cmd:"" description:"Run the server locally"`
	Version     kong.VersionFlag
}

func main() {
	kCtx := kong.Parse(&cliOptions, kong.Vars{
		"version": version.Version,
	})
	kCtx.FatalIfErrorf(kCtx.Run())
}
