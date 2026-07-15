package quiz_api

import (
	"github.com/ed-evo/hankinson/server/ent"
	"github.com/go-chi/chi/v5"
)

var BasePath string = "/quiz"

func QuizGroup(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Mount("/capitoli", getCapitoliRouter(db))
	r.Mount("/domande", getDomandeRouter(db))
	r.Mount("/argomenti", getArgomentiRouter(db))
	r.Mount("/prove", getProveRouter(db))
	return r
}
