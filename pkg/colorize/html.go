package colorize

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type htmlColorizer struct {
	funcName string
}

func NewHTMLColorizer(opts ...ColorizerOption) Colorizer {
	c := &htmlColorizer{
		funcName: "colorize",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (h *htmlColorizer) setCustomFunction(name string) {
	h.funcName = name
}

func (h *htmlColorizer) Format(str string) string {
	t, err := h.newTemplate(str)
	if err != nil {
		fmt.Println(err)
		return "!INVALID"
	}

	b := &bytes.Buffer{}
	if err = t.Execute(b, nil); err != nil {
		fmt.Println(err)
		return "!INVALID"
	}

	return b.String()
}

func (h *htmlColorizer) newTemplate(str string) (*template.Template, error) {
	return template.New("htmlColor").Funcs(template.FuncMap{
		h.funcName: h.colorize,
	}).Parse(str)
}

func (h *htmlColorizer) colorize(format string, text string) string {
	return fmt.Sprintf("<span class=%q>%s</span>", strings.ReplaceAll(format, ";", " "), text)
}
