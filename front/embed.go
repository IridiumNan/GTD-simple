package front

import (
	"embed"
)

//go:embed index.html js style
var FS embed.FS
