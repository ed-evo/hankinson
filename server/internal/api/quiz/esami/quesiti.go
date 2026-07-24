package esami_api

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/attivitaquesitoesame"
	"github.com/ed-evo/hankinson/server/ent/domanda"
	"github.com/ed-evo/hankinson/server/ent/quesitoesame"
	api_utils "github.com/ed-evo/hankinson/server/internal/api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type QuesitoAttivitaBody struct {
	Tipo         string    `json:"tipo"`
	RispostaData *bool     `json:"risposta_data,omitempty"`
	Inizio       time.Time `json:"inizio"`
	DurataMS     int       `json:"durata_ms"`
}

func (q *QuesitoAttivitaBody) Bind(r *http.Request) error {
	return nil
}

func equalBoolPtr(a, b *bool) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return *a == *b
}

func newQuesitiRouter(db *ent.Client) chi.Router {
	quesitiRouter := chi.NewRouter()

	quesitiRouter.Get("/stats", getStats(db))
	quesitiRouter.Route("/{quesitoID:[0-9]+}", func(r chi.Router) {
		r.Use(getQuesitoCtx(db))
		r.Put("/attivita", putAttivitaQuesito(db))
	})

	return quesitiRouter
}

func getQuesitoCtx(db *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var q *ent.QuesitoEsame
			quesitoID := chi.URLParam(r, "quesitoID")
			log.Print("hello")
			if quesitoID == "" {
				render.Render(w, r, api_utils.ErrNotFound)
				return
			}

			id, err := strconv.Atoi(quesitoID)
			if err != nil {
				render.Render(w, r, api_utils.ErrInvalidRequest(err))
				return
			}
			q, err = db.QuesitoEsame.Get(r.Context(), id)
			if ent.IsNotFound(err) {
				render.Render(w, r, api_utils.ErrNotFound)
				return
			}
			if err != nil {
				render.Render(w, r, api_utils.ErrInvalidRequest(err))
				return
			}

			ctx := context.WithValue(r.Context(), "quesito", q)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func putAttivitaQuesito(db *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.Context().Value("quesito").(*ent.QuesitoEsame)
		body := &QuesitoAttivitaBody{}
		if err := render.Bind(r, body); err != nil {
			render.Render(w, r, api_utils.ErrInvalidRequest(err))
			return
		}

		rispostaData := body.RispostaData

		if !equalBoolPtr(rispostaData, q.RispostaFinale) {
			var err error
			q, err = db.QuesitoEsame.UpdateOneID(q.ID).
				SetNillableRispostaFinale(rispostaData).
				Save(r.Context())
			if err != nil {
				render.Render(w, r, api_utils.ErrInternal(err))
				return
			}
		}

		a, err := db.AttivitaQuesitoEsame.Create().
			SetQuesitoEsameID(q.ID).
			SetTipo(attivitaquesitoesame.Tipo(body.Tipo)).
			SetInizio(body.Inizio).
			SetDurataMs(body.DurataMS).
			SetNillableRispostaData(rispostaData).
			Save(r.Context())
		if err != nil {
			render.Render(w, r, api_utils.ErrInternal(err))
			return
		}
		log.Printf("Attivita: %v -> Quesito %v: %v", a.ID, q.ID, body.Tipo)

	}
}

type QuesitoStatsResponse struct {
	Totale    int `json:"totale"`
	Corrette  int `json:"corrette"`
	Sbagliate int `json:"sbagliate"`
	NonDate   int `json:"non_date"`
}

func (s *QuesitoStatsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func getStats(db *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		q := sql.Table(quesitoesame.Table).As("q")
		d := sql.Table(domanda.Table).As("d")
		it := d.C(domanda.FieldIsTrue)
		rf := q.C(quesitoesame.FieldRispostaFinale)

		selector := sql.Select(
			sql.As(sql.Count("*"), "totale"),
			sql.As(sql.Count("CASE WHEN "+rf+" = "+it+" THEN 1 END"), "corrette"),
			sql.As(sql.Count("CASE WHEN "+rf+" <> "+it+" THEN 1 END"), "sbagliate"),
			sql.As(sql.Count("CASE WHEN "+rf+" IS NULL THEN 1 END"), "non_date"),
		).
			From(q).
			Join(d).
			On(
				q.C(quesitoesame.DomandaOriginaleColumn),
				d.C(domanda.FieldID),
			)

		query, _ := selector.Query()

		// quesitoesame.DomandaOriginaleColumn
		rows, err := db.QueryContext(r.Context(), query)
		if err != nil {
			render.Render(w, r, api_utils.ErrInternal(err))
			return
		}
		defer rows.Close()

		for rows.Next() {
			var s QuesitoStatsResponse
			rows.Scan(
				&s.Totale,
				&s.Corrette,
				&s.Sbagliate,
				&s.NonDate,
			)
			render.Render(w, r, &s)
			return
		}
	}
}
