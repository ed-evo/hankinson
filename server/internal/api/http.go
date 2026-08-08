package api

import (
	"net/http"
	"sync"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/esame"
	"github.com/ed-evo/hankinson/server/ent/utente"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	quiz_api "github.com/ed-evo/hankinson/server/internal/api/quiz"
	"github.com/ed-evo/hankinson/server/internal/dto"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var BasePath string = "/api/v1"

func NewApiRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(api_middlewares.CorsHeaders)

	r.Get("/me", meHandler(db))

	r.Group(func(authedRouter chi.Router) {
		authedRouter.Use(api_middlewares.TrivialAuth(db))
		authedRouter.Mount(quiz_api.BasePath, quiz_api.NewQuizRouter(db))
	})

	return r
}

var meMutex sync.Mutex

func meHandler(db *ent.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		meMutex.Lock()
		defer meMutex.Unlock()
		user := api_middlewares.ExtractUser(r.Header)
		if user == nil {
			http.Error(w, "401: Utente obbligatorio", http.StatusUnauthorized)
			return
		}
		userEmail := string(*user)

		ctx := r.Context()
		err := orm.WithTx(ctx, db, func(tx *ent.Tx) error {

			_, err := tx.Utente.Get(ctx, userEmail)

			if ent.IsNotFound(err) {
				err = tx.Utente.Create().
					SetID(userEmail).
					OnConflictColumns(utente.FieldID).
					Ignore().
					Exec(ctx)
			}

			if err != nil {
				return err
			}

			_, err = tx.Esame.Query().
				Where(
					esame.TipoEQ(esame.TipoAperto),
					esame.HasUtenteWith(utente.ID(userEmail)),
				).Only(ctx)

			if ent.IsNotFound(err) {
				err = tx.Esame.Create().
					SetTipo(esame.TipoAperto).
					SetNumeroQuesiti(-1).
					SetMaxErrori(-1).
					SetMinutiDisponibili(-1).
					SetUtenteID(userEmail).
					Exec(ctx)
			}

			return err
		})

		if err != nil {
			dto.RenderError(w, r, err)
			return
		}

		render.JSON(w, r, &user)
	}
}
