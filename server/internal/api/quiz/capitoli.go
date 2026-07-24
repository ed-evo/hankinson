package quiz_api

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"entgo.io/ent/dialect/sql"
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/attivitaquesitoesame"
	"github.com/ed-evo/hankinson/server/ent/capitolo"
	"github.com/ed-evo/hankinson/server/ent/domanda"
	"github.com/ed-evo/hankinson/server/ent/quesitoesame"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	esami_api "github.com/ed-evo/hankinson/server/internal/api/quiz/esami"
	api_utils "github.com/ed-evo/hankinson/server/internal/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newCapitoliRouter(db *ent.Client) chi.Router {
	capitoliRouter := chi.NewRouter()
	api_middlewares.AddToGlobal(capitoliRouter)
	capitoliRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		capitoli, _ := db.Capitolo.Query().
			All(r.Context())
		render.JSON(w, r, capitoli)
	})
	capitoliRouter.Get("/stats", getBasicStats(db))
	capitoliRouter.Route("/{capitoloID:[0-9]+}", func(r chi.Router) {
		r.Use(getCapitoloCtx(db))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			c := r.Context().Value("capitolo").(*ent.Capitolo)

			render.JSON(w, r, c)
		})
	})
	return capitoliRouter
}

func getCapitoloCtx(db *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var c *ent.Capitolo

			if argomentoId := chi.URLParam(r, "capitoloID"); argomentoId != "" {
				id, err := strconv.Atoi(argomentoId)
				if err != nil {
					render.Render(w, r, api_utils.ErrInvalidRequest(err))
					return
				}
				c, err = db.Capitolo.Query().
					Where(capitolo.IDEQ(id)).
					WithDomande(func(q *ent.DomandaQuery) {
						q.WithArgomenti() // Carica la relazione molti-a-molti con gli argomenti
					}).
					Only(r.Context())
				if err != nil {
					render.Render(w, r, api_utils.ErrNotFound)
					return
				}
			}

			ctx := context.WithValue(r.Context(), "capitolo", c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type CapitoliStats struct {
	ID int `json:"id"`
	esami_api.QuesitiStatsResponse
	Durate int `json:"durata"`
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

		query, args := selector.Query()

		log.Printf("%v", args)

		rows, err := db.Debug().QueryContext(
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
				&s.Durate,
			)
			stats = append(stats, &s)
		}

		render.RenderList(w, r, stats)

	}
}
