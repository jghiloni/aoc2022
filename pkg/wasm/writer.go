//go:build wasm && js

package wasm

import (
	"io"
	"strings"
	"syscall/js"
	"time"
)

type htmlWriter struct {
	delay time.Duration
}

func NewWriter(delay time.Duration) io.Writer {
	return &htmlWriter{
		delay: delay,
	}
}

func (h *htmlWriter) Write(ih []byte) (int, error) {
	innerHTML := strings.TrimSuffix(string(ih), "\n")
	retLen := len(innerHTML)

	innerHTML = strings.ReplaceAll(innerHTML, "\n", "<br/>")

	message := map[string]any{
		"type":   "output",
		"output": innerHTML,
	}

	js.Global().Call("postMessage", message)
	time.Sleep(h.delay)

	return retLen, nil
}
