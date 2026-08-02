package quiz_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/spiegazione"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
)

func newDomandeRouter(db *ent.Client) chi.Router {
	domandeRouter := chi.NewRouter()
	ctx := api_context.EntityContextHelper[ent.Domanda, *ent.DomandaQuery]{
		ParamName: "domandaId",
		Fetcher:   &orm.DomandaFetcher{DB: db},
	}

	domandeRouter.Get("/", ctx.JsonListHandler(func(r *http.Request, entities []*ent.Domanda) (any, error) {
		return entities, nil
	}))

	domandeRouter.Route("/{domandaId:[0-9]+}", func(r chi.Router) {
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
