package quiz_api

import (
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/go-chi/chi/v5"
)

var BasePath string = "/quiz"

func NewQuizRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Mount("/capitoli", newCapitoliRouter(db))
	r.Mount("/domande", newDomandeRouter(db))
	r.Mount("/argomenti", newArgomentiRouter(db))
	r.Mount("/prove", newProveRouter(db))
	return r
}
