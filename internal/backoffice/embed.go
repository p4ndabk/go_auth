package backoffice

import "embed"

// The UI is embedded in the binary so `cmd/api` stays self-contained —
// the Docker image copies only the compiled binaries, no asset directory.
//
//go:embed all:assets
var assets embed.FS
