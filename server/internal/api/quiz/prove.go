package quiz_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func getProveRouter(db *ent.Client) chi.Router {
	proveRouter := chi.NewRouter()

	proveRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		domande, _ := db.Domanda.Query().Limit(30).All(r.Context())
		render.JSON(w, r, domande)
	})

	return proveRouter
}
