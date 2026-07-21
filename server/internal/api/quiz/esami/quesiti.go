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
	Tipo   string    `json:"tipo"`
	Inizio time.Time `json:"inizio"`
	Fine   time.Time `json:"fine"`
}

func (q *QuesitoAttivitaBody) Bind(r *http.Request) error {
	return nil
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

			a, err := db.AttivitaQuesitoEsame.Create().
				SetQuesitoEsameID(q.ID).
				SetTipo(attivitaquesitoesame.Tipo(body.Tipo)).
				SetInizio(body.Inizio).
				SetFine(body.Fine).
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
