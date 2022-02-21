package embedded

import "embed"

//go:embed sqlite immudb
var Files embed.FS
