package quiz_api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/internal/ai"
	api_utils "github.com/ed-evo/hankinson/server/internal/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newDomandeRouter(db *ent.Client) chi.Router {
	domandeRouter := chi.NewRouter()

	domandeRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		domande, _ := db.Domanda.Query().All(r.Context())
		render.JSON(w, r, domande)
	})

	domandeRouter.Route("/{domandaId:[0-9]+}", func(r chi.Router) {
		r.Use(getDomandaCtx(db))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			d := r.Context().Value("domanda").(*ent.Domanda)
			render.JSON(w, r, d)
		})
		r.Post("/spiegazione", spiegaDomanda)
	})

	return domandeRouter
}

func getDomandaCtx(db *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var d *ent.Domanda
			if domandaId := chi.URLParam(r, "domandaId"); domandaId != "" {
				id, err := strconv.Atoi(domandaId)
				if err != nil {
					render.Render(w, r, api_utils.ErrInternal(err))
				}

				d, err = db.Domanda.Get(r.Context(), id)
				if err != nil {
					render.Render(w, r, api_utils.ErrNotFound)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "domanda", d)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func spiegaDomanda(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	d := ctx.Value("domanda").(*ent.Domanda)
	gem, err := ai.GetGemini(ctx)
	if err != nil {
		render.Render(w, r, api_utils.ErrInternal(err))
		return
	}
	s, err := gem.Spiega(ctx, d)
	if err != nil {
		render.Render(w, r, api_utils.ErrInternal(err))
	}
	render.JSON(w, r, s)
}
