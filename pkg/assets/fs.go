package assets

import "embed"

//go:embed css/*.css css/*.ttf *.html js
var StaticHTML embed.FS
