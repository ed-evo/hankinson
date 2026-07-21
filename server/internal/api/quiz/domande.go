package quiz_api

import (
	"net/http"
	"strconv"

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

	domandeRouter.Get("/{domandaId:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		argomentoId := chi.URLParam(r, "domandaId")
		if argomentoId == "" {
			http.Error(w, "Domanda non trovata", http.StatusNotFound)
			return
		}
		id, _ := strconv.Atoi(argomentoId)
		d, _ := db.Domanda.Get(r.Context(), id)
		render.JSON(w, r, d)
	})

	return domandeRouter
}
