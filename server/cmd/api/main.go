package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ed-evo/hankinson/server/internal/api"
	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/ed-evo/hankinson/server/webui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	db := orm.GetClient(ctx)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer) // Prevents server crashes if code panics
	r.Mount(api.BasePath, api.NewApiRouter(db))
	r.Mount("/", webui.NewWebRouter())

	// 4. Start Server
	log.Println("🚀 Hankinson backend spinning up on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
