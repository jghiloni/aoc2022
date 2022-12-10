package colorize

import (
	"bytes"
	"fmt"
	"text/template"
)

type nopColorizer struct {
	funcName string
}

func NewNoopColorizer(opts ...ColorizerOption) Colorizer {
	c := &nopColorizer{
		funcName: "colorize",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (n *nopColorizer) setCustomFunction(name string) {
	n.funcName = name
}

func (a *nopColorizer) Format(str string) string {
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

func (n *nopColorizer) newTemplate(str string) (*template.Template, error) {
	return template.New("nopColor").Funcs(template.FuncMap{
		n.funcName: n.colorize,
	}).Parse(str)
}

func (n *nopColorizer) colorize(format string, text string) string {
	return text
}
