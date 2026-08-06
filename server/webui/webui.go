package webui

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
)

//go:embed all:public/*
var PublicFiles embed.FS

func NewWebRouter() chi.Router {
	r := chi.NewRouter()

	publicFS, err := fs.Sub(PublicFiles, "public")
	if err != nil {
		log.Fatalf("FATAL: Cannot access frontend files: %v", err)
	}
	fileServer := http.FileServer(http.FS(publicFS))

	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use a response writer recorder to capture the fileServer status code
		rec := &statusRecorder{ResponseWriter: w, status: 200}
		fileServer.ServeHTTP(rec, r)

		// If http.FileServer returned a 404 (file not found), fallback to index.html
		if rec.status == http.StatusNotFound {
			indexBytes, err := fs.ReadFile(publicFS, "index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(indexBytes)
		}
	}))

	return r
}
// Simple wrapper to record the HTTP status code written by http.FileServer
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	if status != http.StatusNotFound {
		r.ResponseWriter.WriteHeader(status)
	}
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == http.StatusNotFound {
		return len(b), nil // Suppress the default "404 page not found" text output
	}
	return r.ResponseWriter.Write(b)
}

func GetQuizImage(h string) ([]byte, error) {
	p := path.Join("public", "quiz_assets", h+".png")
	return PublicFiles.ReadFile(p)
}
