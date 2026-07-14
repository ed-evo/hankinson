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

	r.Use(middleware.Logger)    // Logs all incoming requests
	r.Use(middleware.Recoverer) // Prevents server crashes if code panics

	r.Mount(api.BasePath, api.ApiGroup(db))

	publicFS, err := webui.Get()
	if err != nil {
		log.Fatalf("FATAL: Cannot access frontend files: %v", err)
	}
	// Direct any non-API web traffic to your frontend files
	r.Handle("/*", http.FileServer(http.FS(publicFS)))

	// 4. Start Server
	log.Println("🚀 Hankinson backend spinning up on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
