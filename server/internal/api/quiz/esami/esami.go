package esami_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/esame"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
)

func newEsamiRouter(db *ent.Client) chi.Router {
	esamiRouter := chi.NewRouter()

	esamiRouter.Route("/{esameID:[0-9]+}", func(r chi.Router) {
		ctx := api_context.EntityContextHelper[ent.Esame, *ent.EsameQuery]{
			ParamName: "esameID",
			Fetcher:   &orm.EsameFetcher{DB: db},
		}

		r.Get("/", ctx.JsonHandler(func(r *http.Request, entity *ent.Esame) (any, error) {
			return entity, nil
		}))

		r.Get("/quesiti", ctx.JsonHandler(func(r *http.Request, entity *ent.Esame) (any, error) {
			return db.Esame.Query().Where(esame.ID(entity.ID)).
				QueryQuesiti().
				WithDomandaOriginale().
				All(r.Context())
		}))
	})

	return esamiRouter
}
