package quiz_api

import (
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	api_middlewares "github.com/ed-evo/hankinson/server/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var BasePath string = "/quiz"

func NewQuizRouter(db *ent.Client) chi.Router {
	r := chi.NewRouter()
	r.Use(api_middlewares.TrivialAuth)
	r.Get("/me", meHandler)
	r.Mount("/capitoli", newCapitoliRouter(db))
	r.Mount("/domande", newDomandeRouter(db))
	r.Mount("/argomenti", newArgomentiRouter(db))
	r.Mount("/prove", newProveRouter(db))
	return r
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	user := api_middlewares.GetUser(r)
	render.JSON(w, r, &user)
}
