package esami_api

import (
	"context"
	"iter"
	"log"
	"maps"
	"math/rand/v2"
	"net/http"
	"slices"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/esame"
	"github.com/ed-evo/hankinson/server/ent/utente"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func newApertoRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Post("/next", next(db))
	return r
}

type QuesitoApertoBody struct {
	Capitoli []int `json:"capitoli"`
}

func (q *QuesitoApertoBody) Bind(r *http.Request) error {
	return nil
}

func selectRandomFromList[T any](items []T) (T, bool) {
	var zero T
	if len(items) == 0 {
		return zero, false
	}
	return items[rand.N(len(items))], true
}

func selectRandomFromSeq[T any](s iter.Seq[T]) (T, bool) {
	return selectRandomFromList(slices.Collect(s))
}

func selectRandoMapKey[K comparable, V any](m map[K]V) (K, bool) {
	return selectRandomFromSeq(maps.Keys(m))
}

func next(db *ent.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := api_middlewares.GetUser(r)
		if user == nil {
			http.Error(w, "Utente Non trovato", http.StatusBadRequest)
		}

		ctx := r.Context()

		esameAperto, err := getArgomentoAperto(ctx, db, string(*user))

		if ent.IsNotFound(err) {
			http.Error(w, "Esame aperto non trovato", http.StatusNotFound)
			return
		}

		body := &QuesitoApertoBody{}

		if err := render.Bind(r, body); err != nil {
			http.Error(w, "Error parsing body", http.StatusBadRequest)
			return
		}

		var capitolo int
		var ok bool

		if body.Capitoli == nil {
			capitolo, ok = selectRandoMapKey(orm.DomandeByCapitolo)
		} else {
			capitolo, ok = selectRandomFromList(body.Capitoli)
		}

		if !ok {
			http.Error(w, "Errore selezione capitolo", http.StatusInternalServerError)
			return
		}

		domandeIDs, ok := orm.DomandeByCapitolo[capitolo]
		if !ok {
			http.Error(w, "Errore recupero domande Ids", http.StatusInternalServerError)
			return
		}

		domandaID, ok := selectRandomFromList(domandeIDs)
		if !ok {
			http.Error(w, "Errore selezione domanda", http.StatusInternalServerError)
			return
		}

		quesito, err := db.QuesitoEsame.Create().
			SetEsameID(esameAperto.ID).
			SetDomandaOriginaleID(domandaID).
			Save(ctx)

		if err != nil {
			http.Error(w, "Errore Creazione Quesito", http.StatusInternalServerError)
			return
		}

		log.Printf("Esame Aperto: %v Capitolo %v -> Domanda: %v", esameAperto.ID, capitolo, domandaID)

		response := struct {
			ID        int `json:"id"`
			EsameId   int `json:"esameId"`
			DomandaId int `json:"domandaId"`
		}{
			ID:        quesito.ID,
			EsameId:   esameAperto.ID,
			DomandaId: domandaID,
		}
		render.JSON(w, r, response)

	}
}

func getArgomentoAperto(ctx context.Context, db *ent.Client, userEmail string) (*ent.Esame, error) {
	return db.Esame.Query().
		Where(
			esame.TipoEQ(esame.TipoAperto),
			esame.HasUtenteWith(utente.IDEQ(userEmail)),
		).First(ctx)
}
