package quiz_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newArgomentiRouter(db *ent.Client) chi.Router {
	argomentiRouter := chi.NewRouter()

	argomentiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		argomenti, _ := db.Argomento.Query().All(r.Context())
		render.JSON(w, r, argomenti)
	})

	argomentiRouter.Route("/{argomentoID:[0-9]+}", func(r chi.Router) {
		ctx := api_context.EntityContextHelper[ent.Argomento]{
			ParamName:  "argomentoID",
			ContextKey: "argomento",
			Fetcher:    orm.ArgomentoFetcher{DB: db},
		}
		r.Use(ctx.Middleware())
		r.Get("/", ctx.JsonHandler(func(r *http.Request, entity *ent.Argomento) (any, error) {
			return entity, nil
		}))
	})

	return argomentiRouter
}
