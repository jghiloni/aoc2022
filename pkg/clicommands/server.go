package clicommands

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alecthomas/kong"
	"github.com/jghiloni/aoc2022/pkg/assets"
)

type LocalServerCommand struct {
	Port int `short:"p" default:"3000"`
}

func (l *LocalServerCommand) Run(kCtx *kong.Context) error {
	log.Println("listening on port", l.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", l.Port), http.FileServer(http.FS(assets.StaticHTML)))
}
