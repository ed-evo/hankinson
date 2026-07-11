package webui

import (
	"embed"
	"io/fs"
)

//go:embed public/*
var PublicFiles embed.FS

var publicFS, err = fs.Sub(PublicFiles, "public")

func Get() (fs.FS, error) {
	return publicFS, err
}
