package esami_api

import (
	"context"
	"log"
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/esame"
	"github.com/ed-evo/hankinson/server/ent/quesitoesame"
	"github.com/ed-evo/hankinson/server/ent/utente"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	"github.com/ed-evo/hankinson/server/internal/dto"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newApertoRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Post("/next", next(db))
	return r
}

type EsameApertoBody struct {
	Capitoli []int `json:"capitoli"`
}

func (q *EsameApertoBody) Bind(r *http.Request) error {
	return nil
}

func next(db *ent.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := api_middlewares.GetUser(r)
		if user == nil {
			http.Error(w, "Utente Non trovato", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		esameAperto, err := getEsameAperto(ctx, db, string(*user))

		if ent.IsNotFound(err) {
			http.Error(w, "Esame aperto non trovato", http.StatusNotFound)
			return
		}

		body := &EsameApertoBody{}

		if err := render.Bind(r, body); err != nil {
			dto.RenderError(w, r, err)
			return
		}

		domandaIDs, err := orm.RandomDomandeIds(body.Capitoli, 1)
		if err != nil {
			dto.RenderError(w, r, err)
			return
		}

		domandaID := domandaIDs[0]

		quesito, err := db.QuesitoEsame.Create().
			SetEsameID(esameAperto.ID).
			SetDomandaOriginaleID(domandaID).
			Save(ctx)

		if err != nil {
			dto.RenderError(w, r, err)
			return
		}

		log.Printf("Esame Aperto: %v -> Domanda: %v", esameAperto.ID, domandaID)

		response, err := db.QuesitoEsame.Query().
			Where(quesitoesame.ID(quesito.ID)).
			Only(ctx)
		if err != nil {
			dto.RenderError(w, r, err)
			return 
		}
		render.JSON(w, r, response)

	}
}

func getEsameAperto(ctx context.Context, db *ent.Client, userEmail string) (*ent.Esame, error) {
	return db.Esame.Query().
		Where(
			esame.TipoEQ(esame.TipoAperto),
			esame.HasUtenteWith(utente.IDEQ(userEmail)),
		).First(ctx)
}
