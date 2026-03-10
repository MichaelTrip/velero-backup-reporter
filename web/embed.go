package web

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var distFiles embed.FS

// DistFS returns the embedded SPA distribution filesystem.
var DistFS, _ = fs.Sub(distFiles, "dist")
