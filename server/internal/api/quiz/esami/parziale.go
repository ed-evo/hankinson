package esami_api

import (
	"log"
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/esame"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/ed-evo/hankinson/server/internal/utils"
	"github.com/ed-evo/hankinson/server/pkg/api/api_errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newParzialiRouter(db *ent.Client) chi.Router {
	parzialiRouter := chi.NewRouter()

	parzialiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		all, err := db.Esame.Query().
			Where(esame.TipoEQ(esame.TipoParziale)).
			All(r.Context())
		if err != nil {
			render.Render(w, r, api_errors.ErrInternal(err))
			return
		}
		render.JSON(w, r, all)
	})
	parzialiRouter.Put("/", createParziale(db))

	return parzialiRouter
}

type EsameParzialeBody struct {
	Capitoli          []int `json:"capitoli"`
	NumeroQuesiti     int   `json:"numero_quesiti"`
	MaxErrori         int   `json:"max_errori"`
	MinutiDisponibili int   `json:"minuti_disponibili"`
}

func (p *EsameParzialeBody) Bind(r *http.Request) error {
	p.NumeroQuesiti = utils.Clamp(p.NumeroQuesiti, 5, 50)
	p.MaxErrori = utils.Clamp(p.MaxErrori, 0, min(10, p.NumeroQuesiti/3))
	p.MinutiDisponibili = utils.Clamp(p.MinutiDisponibili, 5, p.NumeroQuesiti)
	return nil
}

func createParziale(db *ent.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := &EsameParzialeBody{}
		if err := render.Bind(r, body); err != nil {
			render.Render(w, r, api_errors.ErrInvalidRequest(err))
			return
		}
		log.Printf("Request %v", body)

		domandeIDs, err := orm.RandomDomandeIds(body.Capitoli, body.NumeroQuesiti)
		if err != nil {
			render.Render(w, r, api_errors.ErrInternal(err))
			return
		}

		log.Printf("Domande IDs: %v", domandeIDs)

		user := api_middlewares.GetUser(r)

		ctx := r.Context()

		var e *ent.Esame

		type QuesitoItem struct {
			ID        int `json:"id"`
			EsameId   int `json:"esameId"`
			DomandaId int `json:"domandaId"`
		}

		err = orm.WithTx(ctx, db, func(tx *ent.Tx) error {
			e, err = tx.Esame.Create().
				SetTipo(esame.TipoParziale).
				SetNumeroQuesiti(body.NumeroQuesiti).
				SetUtenteID(string(*user)).
				SetMaxErrori(body.MaxErrori).
				SetMinutiDisponibili(body.MinutiDisponibili).
				Save(ctx)
			if err != nil {
				return err
			}

			builders := make([]*ent.QuesitoEsameCreate, len(domandeIDs))
			for i, domandaID := range domandeIDs {
				builders[i] = tx.QuesitoEsame.Create().
					SetEsameID(e.ID).
					SetDomandaOriginaleID(domandaID)
			}

			err := tx.QuesitoEsame.CreateBulk(builders...).Exec(ctx)
			if err != nil {
				return err
			}

			if err := tx.Commit(); err != nil {
				return err
			}
			return nil
		})

		render.JSON(w, r, e)
	}
}
