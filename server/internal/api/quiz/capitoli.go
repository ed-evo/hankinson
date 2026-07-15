package quiz_api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/capitolo"
	api_utils "github.com/ed-evo/hankinson/server/internal/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func getCapitoliRouter(db *ent.Client) chi.Router {
	capitoliRouter := chi.NewRouter()
	capitoliRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		capitoli, _ := db.Capitolo.Query().
			All(r.Context())
		render.JSON(w, r, capitoli)
	})
	capitoliRouter.Route("/{capitoloID:[0-9]+}", func(r chi.Router) {
		r.Use(getCapitoloCtx(db))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			c := r.Context().Value("capitolo").(*ent.Capitolo)

			render.JSON(w, r, c)
		})
	})
	return capitoliRouter
}

func getCapitoloCtx(db *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var c *ent.Capitolo

			if argomentoId := chi.URLParam(r, "capitoloID"); argomentoId != "" {
				id, err := strconv.Atoi(argomentoId)
				if err != nil {
					render.Render(w, r, api_utils.ErrInvalidRequest(err))
					return
				}
				c, err = db.Capitolo.Query().
					Where(capitolo.IDEQ(id)).
					WithDomande(func(q *ent.DomandaQuery) {
						q.WithArgomenti() // Carica la relazione molti-a-molti con gli argomenti
					}).
					Only(r.Context())
				if err != nil {
					render.Render(w, r, api_utils.ErrNotFound)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "capitolo", c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
