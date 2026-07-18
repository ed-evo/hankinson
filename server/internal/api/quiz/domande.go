package quiz_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newDomandeRouter(db *ent.Client) chi.Router {
	domandeRouter := chi.NewRouter()

	domandeRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		domande, _ := db.Domanda.Query().All(r.Context())
		render.JSON(w, r, domande)
	})

	return domandeRouter
}
