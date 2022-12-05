//go:build wasm && js

package wasm

import (
	"fmt"
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

	parser := js.Global().Get("DOMParser").New()
	parsed := parser.Call("parseFromString", innerHTML, "text/html")
	children := parsed.Call("getElementsByTagName", "body").Index(0).Get("childNodes")

	fragment := h.document.Call("createDocumentFragment")
	for children.Length() > 0 {
		fragment.Call("append", children.Index(0))
	}

	fmt.Println(fragment)

	div := h.document.Call("createElement", "div")
	div.Set("className", "line")
	div.Call("appendChild", fragment)

	h.console.Call("appendChild", div)
	return len(ih), nil
}
