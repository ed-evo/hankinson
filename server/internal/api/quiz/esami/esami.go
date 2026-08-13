package esami_api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/correzione"
	"github.com/ed-evo/hankinson/server/ent/esame"
	"github.com/ed-evo/hankinson/server/internal/ai"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	"github.com/ed-evo/hankinson/server/internal/dto"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newEsamiRouter(db *ent.Client) chi.Router {
	esamiRouter := chi.NewRouter()

	esamiRouter.Route("/{esameID:[0-9]+}", func(r chi.Router) {
		ctx := api_context.EntityContextHelper[ent.Esame, *ent.EsameQuery]{
			ParamName: "esameID",
			Fetcher:   &orm.EsameFetcher{DB: db},
		}

		r.Get("/", ctx.JsonHandler(
			func(r *http.Request, entity *ent.Esame) (any, error) {
				return entity, nil
			},
		))

		r.Get("/quesiti", ctx.JsonHandler(func(r *http.Request, entity *ent.Esame) (any, error) {
			return db.Esame.Query().Where(esame.ID(entity.ID)).
				QueryQuesiti().
				All(r.Context())
		}))

		r.Get(
			"/correzioni",
			ctx.JsonHandler(
				func(r *http.Request, entity *ent.Esame) (any, error) {
					return entity.Edges.Correzioni, nil
				},
				func(q *ent.EsameQuery) *ent.EsameQuery {
					return q.WithCorrezioni()
				},
			),
		)

		r.With(middleware.Timeout(10 * time.Minute)).Post(
			"/ai-corregge",
			ctx.JsonHandler(
				func(r *http.Request, entity *ent.Esame) (any, error) {
					if entity.Tipo == esame.TipoAperto {
						return nil, dto.ErrInvalidRequest(fmt.Errorf("Esami di tipo Aperto non supportano la correzione."))
					}
					if len(entity.Edges.Correzioni) > 0 {
						return entity.Edges.Correzioni, nil
					}
					ctx := r.Context()
					gem, err := ai.GetGemini(ctx)
					if err != nil {
						return nil, err
					}
					c, err := gem.Correggi(ctx, db, entity.ID)
					if err != nil {
						return nil, err
					}
					newC, err := db.Correzione.Create().
						SetEsameID(entity.ID).
						SetType(correzione.TypeAi).
						SetEsaminatore(c.Model).
						SetIsPromosso(c.Promosso).
						SetMeta(c.Metadata).
						SetTesto(c.Testo).
						Save(ctx)
					if err != nil {
						return nil, err
					}
					return []*ent.Correzione{ newC }, nil

				},
				func(q *ent.EsameQuery) *ent.EsameQuery {
					return q.WithCorrezioni()
				},
			),
		)
	})

	return esamiRouter
}
