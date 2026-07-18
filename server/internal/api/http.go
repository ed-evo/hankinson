package api

import (
	"github.com/ed-evo/hankinson/server/ent"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	quiz_api "github.com/ed-evo/hankinson/server/internal/api/quiz"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var BasePath string = "/api/v1"

func NewApiRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(api_middlewares.CorsHeaders)
	r.Mount(quiz_api.BasePath, quiz_api.NewQuizRouter(db))

	return r
}
