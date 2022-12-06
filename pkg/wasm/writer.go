//go:build wasm && js

package wasm

import (
	"io"
	"strings"
	"syscall/js"
)

type htmlWriter struct {
	console  js.Value
	document js.Value
}

func NewWriter(console js.Value) io.Writer {
	return &htmlWriter{
		console:  console,
		document: js.Global().Get("document"),
	}
}

func (h *htmlWriter) Write(ih []byte) (int, error) {
	innerHTML := strings.TrimSuffix(string(ih), "\n")
	innerHTML = strings.ReplaceAll(innerHTML, "\n", "<br/>")

	div := h.document.Call("createElement", "div")
	div.Set("className", "line")
	div.Set("innerHTML", innerHTML)

	eventOptions := map[string]any{
		"detail": div,
	}

	event := js.Global().Get("CustomEvent").New("output", eventOptions)
	h.console.Call("dispatchEvent", event)

	return len(ih), nil
}
