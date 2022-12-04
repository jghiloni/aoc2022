//go:build wasm && js

package utils

import (
	"errors"
	"io"
	"syscall/js"
)

type htmlWriter struct {
	console js.Value
}

func NewWriter(console js.Value) io.Writer {
	return &htmlWriter{
		console: console,
	}
}

func (h *htmlWriter) Write(b []byte) (int, error) {
	innerHTML := h.console.Get("innerHTML")
	if !innerHTML.Truthy() {
		return 0, errors.New("could not get innerHTML attribute")
	}

	html := innerHTML.String() + ANSItoHTML(b)
	h.console.Set("innerHTML", html)

	return len(html), nil
}
