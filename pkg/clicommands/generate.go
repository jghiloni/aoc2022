package clicommands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alecthomas/kong"
)

type GenerateExerciseCommand struct {
	Day            int    `short:"d" description:"The day to generate the exercise for"`
	RepositoryRoot string `short:"o" description:"The repository root" default:"."`
	SessionCookie  string `short:"c" env:"AOC_SESSION_COOKIE"`
}

func (g *GenerateExerciseCommand) Run(kCtx *kong.Context) error {
	if g.Day <= 0 {
		t := time.Now()
		g.Day = t.Day()
	}

	fmt.Fprintf(kCtx.Stdout, "Generating code stub for day %d\n", g.Day)

	tmpl, err := template.ParseFiles(filepath.Join(g.RepositoryRoot, "exercise.template"))
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	outFile, err := os.Create(filepath.Join(g.RepositoryRoot, "pkg", "exercise", fmt.Sprintf("day%d.go", g.Day)))
	if err != nil {
		return fmt.Errorf("could not create ouptut file: %w", err)
	}
	defer outFile.Close()

	if err := tmpl.ExecuteTemplate(outFile, "exercise.template", g); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	if g.SessionCookie == "" {
		fmt.Fprintf(kCtx.Stdout, "Get your input from https://adventofcode.com/2022/day/%[1]d/input and put it in %[2]s/pkg/inputs/day%[1]d.txt", g.Day, g.RepositoryRoot)
		return nil
	}

	fmt.Fprintln(kCtx.Stdout, "Fetching input text")
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://adventofcode.com/2022/day/%d/input", g.Day), nil)
	if err != nil {
		return fmt.Errorf("could not build input request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: g.SessionCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(kCtx.Stdout, "Get your input from https://adventofcode.com/2022/day/%[1]d/input and put it in %[2]s/pkg/inputs/day%[1]d.txt", g.Day, g.RepositoryRoot)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(kCtx.Stdout, "Get your input from https://adventofcode.com/2022/day/%[1]d/input and put it in %[2]s/pkg/inputs/day%[1]d.txt", g.Day, g.RepositoryRoot)
		return nil
	}

	inputFile, err := os.Create(filepath.Join(g.RepositoryRoot, "pkg", "inputs", fmt.Sprintf("day%d.txt", g.Day)))
	if err != nil {
		fmt.Fprintf(kCtx.Stdout, "Get your input from https://adventofcode.com/2022/day/%[1]d/input and put it in %[2]s/pkg/inputs/day%[1]d.txt", g.Day, g.RepositoryRoot)
		return nil
	}
	defer inputFile.Close()

	_, err = io.Copy(inputFile, resp.Body)
	return err
}
