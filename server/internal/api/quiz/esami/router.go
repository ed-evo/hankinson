package esami_api

import (
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/go-chi/chi/v5"
)

func NewEsamiRouter(db *ent.Client) chi.Router {
	esameRouter := chi.NewRouter()

	esameRouter.Mount("/", newEsamiRouter(db))

	esameRouter.Mount("/aperto", newApertoRouter(db))

	esameRouter.Mount("/parziali", newParzialiRouter(db))

	esameRouter.Mount("/quesiti", newQuesitiRouter(db))

	return esameRouter
}
