package api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	quiz_api "github.com/ed-evo/hankinson/server/internal/api/quiz"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var BasePath string = "/api/v1"

// 2. Simple CORS middleware for local Vite development
func CorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func NewApiRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(CorsHeaders)
	r.Mount(quiz_api.BasePath, quiz_api.NewQuizRouter(db))

	return r
}
