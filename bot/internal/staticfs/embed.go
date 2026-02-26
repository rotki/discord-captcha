package staticfs

import "embed"

//go:embed all:files
var StaticFiles embed.FS
