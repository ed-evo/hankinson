package quiz_api

import (
	"fmt"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/attivitaquesitoesame"
	"github.com/ed-evo/hankinson/server/ent/domanda"
	"github.com/ed-evo/hankinson/server/ent/esame"
	"github.com/ed-evo/hankinson/server/ent/quesitoesame"
	api_context "github.com/ed-evo/hankinson/server/internal/api/context"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	esami_api "github.com/ed-evo/hankinson/server/internal/api/quiz/esami"
	"github.com/ed-evo/hankinson/server/internal/dto"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newCapitoliRouter(db *ent.Client) chi.Router {
	capitoliRouter := chi.NewRouter()
	ctx := api_context.EntityContextHelper[ent.Capitolo, *ent.CapitoloQuery]{
		ParamName: "capitoloID",
		Fetcher:   &orm.CapitoloFetcher{DB: db},
	}

	capitoliRouter.Get("/stats", getBasicStats(db))
	// chaced responses
	capitoliRouter.Group(func(r chi.Router) {
		api_middlewares.AddToGlobal(r)
		r.Get("/", ctx.JsonListHandler(func(r *http.Request, entities []*ent.Capitolo) (any, error) {
			return entities, nil
		}))
		r.Route("/{capitoloID:[0-9]+}", func(r chi.Router) {
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

		u := api_middlewares.GetUser(r)
		if u == nil {
			dto.RenderError(w, r, dto.ErrInvalidRequest(fmt.Errorf("User required")))
			return
		}
		e := sql.Table(esame.Table).As("e")
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
			).
			AppendSelectExpr(sql.ExprP(
				sql.As(
					sql.Sum("CASE WHEN "+aT+" = ? THEN "+aD+" ELSE 0 END")+" + "+
					sql.Sum("CASE WHEN "+aT+" = ? THEN "+aD+" ELSE 0 END"),
					"tempo",
				),
				attivitaquesitoesame.TipoSalta,
				attivitaquesitoesame.TipoRisposta,
			)).
			From(d).
			Join(q).On(q.C(quesitoesame.DomandaOriginaleColumn), dID).
			Join(e).On(q.C(quesitoesame.EsameColumn), e.C(esame.FieldID)).
			LeftJoin(a).On(q.C(quesitoesame.FieldID), a.C(attivitaquesitoesame.QuesitoEsameColumn)).
			Where(sql.EQ(e.C(esame.UtenteColumn), *u),).
			GroupBy(cID)

		query, args := selector.Query()

		rows, err := db.QueryContext(
			r.Context(),
			query,
			args...,
		)
		if err != nil {
			dto.RenderError(w, r, err)
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
