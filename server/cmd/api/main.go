package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ed-evo/hankinson/server/internal/orm"
	"github.com/ed-evo/hankinson/server/webui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type MessageResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// Handler function for /api/v1/hello
func handleHello(w http.ResponseWriter, r *http.Request) {
	response := MessageResponse{
		Message: "Welcome to the Hankinson App API!",
		Status:  http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the struct directly to the response writer
	json.NewEncoder(w).Encode(response)
}

func main() {
	ctx := context.Background()

	_, err := orm.GetClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create orm client %v", err)
	}

	r := chi.NewRouter()

	// 1. Global Middleware
	r.Use(middleware.Logger)    // Logs all incoming requests
	r.Use(middleware.Recoverer) // Prevents server crashes if code panics

	// 2. Simple CORS middleware for local Vite development
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// 3. Route Grouping for API Versioning
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/hello", handleHello)
	})

	publicFS, err := webui.Get()
	if err != nil {
		log.Fatalf("FATAL: Cannot access frontend files: %v", err)
	}
	// Direct any non-API web traffic to your frontend files
	r.Handle("/*", http.FileServer(http.FS(publicFS)))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Custom 404"))
	})

	// 4. Start Server
	log.Println("🚀 Hankinson backend spinning up on port :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
