package colorize

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

const (
	black int = iota
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

const (
	foreground    = 30
	background    = 40
	bright        = 60
	bold          = 1
	faint         = 2
	italic        = 3
	underline     = 4
	strikethrough = 9
)

type ansiColorizer struct {
	funcName string
}

var codeMap = map[string]int{
	"bold":              bold,
	"faint":             faint,
	"italic":            italic,
	"underline":         underline,
	"strikethrough":     strikethrough,
	"black":             foreground + black,
	"red":               foreground + red,
	"green":             foreground + green,
	"yellow":            foreground + yellow,
	"blue":              foreground + blue,
	"magenta":           foreground + magenta,
	"cyan":              foreground + cyan,
	"white":             foreground + white,
	"bg-black":          background + black,
	"bg-red":            background + red,
	"bg-green":          background + green,
	"bg-yellow":         background + yellow,
	"bg-blue":           background + blue,
	"bg-magenta":        background + magenta,
	"bg-cyan":           background + cyan,
	"bg-white":          background + white,
	"bright-black":      foreground + black + bright,
	"bright-red":        foreground + red + bright,
	"bright-green":      foreground + green + bright,
	"bright-yellow":     foreground + yellow + bright,
	"bright-blue":       foreground + blue + bright,
	"bright-magenta":    foreground + magenta + bright,
	"bright-cyan":       foreground + cyan + bright,
	"bright-white":      foreground + white + bright,
	"bright-bg-black":   background + black + bright,
	"bright-bg-red":     background + red + bright,
	"bright-bg-green":   background + green + bright,
	"bright-bg-yellow":  background + yellow + bright,
	"bright-bg-blue":    background + blue + bright,
	"bright-bg-magenta": background + magenta + bright,
	"bright-bg-cyan":    background + cyan + bright,
	"bright-bg-white":   background + white + bright,
}

func NewANSIColorizer(opts ...ColorizerOption) Colorizer {
	c := &ansiColorizer{
		funcName: "colorize",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (a *ansiColorizer) setCustomFunction(name string) {
	a.funcName = name
}

func (a *ansiColorizer) Format(str string) string {
	t, err := a.newTemplate(str)
	if err != nil {
		return fmt.Sprintln(err)
	}

	b := &bytes.Buffer{}
	if err = t.Execute(b, nil); err != nil {
		return fmt.Sprintln(err)
	}

	return b.String()
}

func (a *ansiColorizer) newTemplate(str string) (*template.Template, error) {
	return template.New("ansiColor").Funcs(template.FuncMap{
		a.funcName: a.colorize,
	}).Parse(str)
}

func (a *ansiColorizer) colorize(format string, text string) string {
	codes := strings.Split(format, ";")
	ansiCodes := make([]string, len(codes))
	for i := range codes {
		ansiCodes[i] = strconv.Itoa(codeMap[codes[i]])
	}

	return fmt.Sprintf("\x1b[%sm%s\x1b[m", strings.Join(ansiCodes, ";"), text)
}
