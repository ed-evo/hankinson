package webui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/internal/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed public/*
var PublicFiles embed.FS

func NewWebRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer) // Prevents server crashes if code panics

	r.Mount(api.BasePath, api.NewApiRouter(db))

	publicFS, err := fs.Sub(PublicFiles, "public")
	if err != nil {
		log.Fatalf("FATAL: Cannot access frontend files: %v", err)
	}
	// Direct any non-API web traffic to your frontend files
	r.Handle("/*", http.FileServer(http.FS(publicFS)))

	return r
}
