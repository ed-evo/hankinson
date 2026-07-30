package quiz_api

import (
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/attivitaquesitoesame"
	"github.com/ed-evo/hankinson/server/ent/domanda"
	"github.com/ed-evo/hankinson/server/ent/quesitoesame"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	esami_api "github.com/ed-evo/hankinson/server/internal/api/quiz/esami"
	api_utils "github.com/ed-evo/hankinson/server/internal/api/utils"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newCapitoliRouter(db *ent.Client) chi.Router {
	capitoliRouter := chi.NewRouter()

	capitoliRouter.Get("/stats", getBasicStats(db))
	// chaced responses
	capitoliRouter.Group(func(r chi.Router) {
		api_middlewares.AddToGlobal(r)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			capitoli, _ := db.Capitolo.Query().
				All(r.Context())
			render.JSON(w, r, capitoli)
		})
		r.Route("/{capitoloID:[0-9]+}", func(r chi.Router) {
			ctx := api_context.EntityContextHelper[ent.Capitolo]{
				ParamName:  "capitoloID",
				ContextKey: "capitolo",
				Fetcher:    orm.CapitoloFetcher{DB: db},
			}
			r.Use(ctx.Middleware())
			r.Get("/", ctx.JsonHandler(func(r *http.Request, entity *ent.Capitolo) (any, error) {
				return entity, nil
			}))
		})
	})
	return capitoliRouter
}

type CapitoliStats struct {
	ID int `json:"id"`
	esami_api.QuesitiStatsResponse
	DurateMs int `json:"durata_ms"`
}

func (s *CapitoliStats) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func getBasicStats(db *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := sql.Table(quesitoesame.Table).As("q")
		qRf := q.C(quesitoesame.FieldRispostaFinale)
		d := sql.Table(domanda.Table).As("d")
		cID := d.C(domanda.CapitoloColumn)
		dID := d.C(domanda.FieldID)
		dIt := d.C(domanda.FieldIsTrue)
		a := sql.Table(attivitaquesitoesame.Table).As("a")
		aD := a.C(attivitaquesitoesame.FieldDurataMs)
		aT := a.C(attivitaquesitoesame.FieldTipo)

		selector := sql.Select(
			cID,
			sql.As(sql.Count(sql.Distinct(dID)), "totale"),
			sql.As(sql.Count(sql.Distinct("CASE WHEN "+qRf+" = "+dIt+" THEN "+dID+" END")), "corrette"),
			sql.As(sql.Count(sql.Distinct("CASE WHEN "+qRf+" <> "+dIt+" THEN "+dID+" END")), "sbagliate"),
			sql.As(sql.Count(sql.Distinct("CASE WHEN "+qRf+" IS NULL THEN "+dID+" END")), "non_date"),
			sql.As(
				sql.Sum("CASE WHEN "+aT+" = $1 THEN "+aD+" ELSE 0 END")+" + "+
					sql.Sum("CASE WHEN "+aT+" = $2 THEN "+aD+" ELSE 0 END")+" - "+
					sql.Sum("CASE WHEN "+aT+" = $3 THEN "+aD+" ELSE 0 END"),
				"tempo",
			),
		).
			From(d).
			Join(q).On(q.C(quesitoesame.DomandaOriginaleColumn), dID).
			LeftJoin(a).On(q.C(quesitoesame.FieldID), a.C(attivitaquesitoesame.QuesitoEsameColumn)).
			GroupBy(cID)

		query, _ := selector.Query()

		rows, err := db.QueryContext(
			r.Context(),
			query,
			attivitaquesitoesame.TipoSalta,
			attivitaquesitoesame.TipoRisposta,
			attivitaquesitoesame.TipoPausa,
		)
		if err != nil {
			render.Render(w, r, api_utils.ErrInternal(err))
			return
		}
		defer rows.Close()

		var stats []render.Renderer

		for rows.Next() {
			var s CapitoliStats
			rows.Scan(
				&s.ID,
				&s.Totale,
				&s.Corrette,
				&s.Sbagliate,
				&s.NonDate,
				&s.DurateMs,
			)
			stats = append(stats, &s)
		}

		render.RenderList(w, r, stats)

	}
}
