package esami_api

import (
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/go-chi/chi/v5"
)

func NewEsamiRouter(db *ent.Client) chi.Router {
	esameRouter := chi.NewRouter()

	esameRouter.Mount("/aperto", newApertoRouter(db))

	return esameRouter
}
