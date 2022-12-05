package colorize

import (
	"errors"
	"io"
)

type Colorizer interface {
	Format(string) string
	setCustomFunction(string)
}

type ColorizerOption func(Colorizer)

func CustomFunctionName(name string) ColorizerOption {
	return func(c Colorizer) {
		c.setCustomFunction(name)
	}
}

type colorWriter struct {
	out io.Writer
	c   Colorizer
}

func NewColorWriter(out io.Writer, c Colorizer) io.Writer {
	return &colorWriter{
		out,
		c,
	}
}

func (c *colorWriter) Write(b []byte) (int, error) {
	if c.out == nil || c.c == nil {
		return 0, errors.New("neither the out writer nor colorizer can be nil")
	}

	s := string(b)
	colorized := c.c.Format(s)

	return c.out.Write([]byte(colorized))
}
