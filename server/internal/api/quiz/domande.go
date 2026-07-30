package quiz_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/spiegazione"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	"github.com/ed-evo/hankinson/server/internal/orm"
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
		ctx := api_context.EntityContextHelper[ent.Domanda]{
			ParamName:  "domandaId",
			ContextKey: "domanda",
			Fetcher:    orm.DomandaFetcher{DB: db},
		}
		r.Use(ctx.Middleware())
		r.Get("/", ctx.JsonHandler(func(r *http.Request, entity *ent.Domanda) (any, error) {
			return entity, nil
		}))
		r.Post("/spiegazione", ctx.JsonHandler(func(r *http.Request, entity *ent.Domanda) (any, error) {
			return db.Spiegazione.Query().
				Where(spiegazione.NumeroDomandaEQ(entity.ID)).
				Only(r.Context())
		}))
	})

	return domandeRouter
}
