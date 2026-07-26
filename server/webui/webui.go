package webui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
)

//go:embed public/*
var PublicFiles embed.FS

func NewWebRouter() chi.Router {
	r := chi.NewRouter()

	publicFS, err := fs.Sub(PublicFiles, "public")
	if err != nil {
		log.Fatalf("FATAL: Cannot access frontend files: %v", err)
	}
	// Direct any non-API web traffic to your frontend files
	r.Handle("/*", http.FileServer(http.FS(publicFS)))

	return r
}

func GetQuizImage(h string) ([]byte, error) {
	p := path.Join("public", "quiz_assets", h+".png")
	return PublicFiles.ReadFile(p)
}
