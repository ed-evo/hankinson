package quiz_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
)

func newArgomentiRouter(db *ent.Client) chi.Router {
	argomentiRouter := chi.NewRouter()
	ctx := api_context.EntityContextHelper[ent.Argomento, *ent.ArgomentoQuery]{
		ParamName: "argomentoID",
		Fetcher:   &orm.ArgomentoFetcher{DB: db},
	}

	argomentiRouter.Get("/", ctx.JsonListHandler(func(r *http.Request, entities []*ent.Argomento) (any, error) {
		return entities, nil
	}))

	argomentiRouter.Route("/{argomentoID:[0-9]+}", func(r chi.Router) {
		r.Get("/", ctx.JsonHandler(func(r *http.Request, entity *ent.Argomento) (any, error) {
			return entity, nil
		}))
	})

	return argomentiRouter
}
