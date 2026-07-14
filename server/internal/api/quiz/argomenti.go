package quiz_api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/argomento"
	api_utils "github.com/ed-evo/hankinson/server/internal/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func getArgomentiRouter(db *ent.Client) chi.Router {
	argomentiRouter := chi.NewRouter()

	argomentiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		argomenti, _ := db.Argomento.Query().All(r.Context())
		render.JSON(w, r, argomenti)
	})

	argomentiRouter.Route("/{argomentoID:[0-9]+}", func(r chi.Router) {
		r.Use(getArgomentoCtx(db))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			a := r.Context().Value("argomento").(*ent.Argomento)

			render.JSON(w, r, a)
		})
	})

	return argomentiRouter
}

func getArgomentoCtx(db *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var a *ent.Argomento

			if argomentoId := chi.URLParam(r, "argomentoID"); argomentoId != "" {
				id, err := strconv.Atoi(argomentoId)
				if err != nil {
					render.Render(w, r, api_utils.ErrInvalidRequest(err))
					return
				}
				a, err = db.Argomento.Query().
					Where(argomento.IDEQ(id)).
					WithDomande().
					Only(r.Context())
				if err != nil {
					render.Render(w, r, api_utils.ErrNotFound)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "argomento", a)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
