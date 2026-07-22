package esami_api

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/attivitaquesitoesame"
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

	quesitiRouter.Route("/{quesitoID:[0-9]+}", func(r chi.Router) {
		r.Use(getQuesitoCtx(db))
		r.Put("/attivita", func(w http.ResponseWriter, r *http.Request) {
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

		})
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
