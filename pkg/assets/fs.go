package assets

import "embed"

//go:embed css/*.css css/*.ttf *.html js/*.js js/aoc.wasm
var StaticHTML embed.FS
