package assets

import "embed"

//go:embed templates
var Templates embed.FS

//go:embed fonts
var Fonts embed.FS

//go:embed images
var Images embed.FS
