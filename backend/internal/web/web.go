package web

import (
	"embed"
	"io/fs"
)

// dist holds the built frontend (populated by the Docker build; the
// checked-in placeholder keeps local `go run` working).
//
//go:embed all:dist
var dist embed.FS

// FS returns the embedded frontend files rooted at the dist directory.
func FS() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		// all:dist guarantees the subdirectory exists
		panic(err)
	}
	return sub
}
