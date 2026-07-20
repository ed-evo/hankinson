package api_middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/ed-evo/hankinson/server/ent"
	"github.com/ed-evo/hankinson/server/ent/utente"
)

type contextKey string

type User string

const UserKey contextKey = "userEmail"

func ExtractUser(h http.Header) *User {
	if h == nil {
		return nil
	}
	userEmail := h.Get("Cf-Access-Authenticated-User-Email")

	if userEmail == "" {
		userEmail = h.Get("X-Authenticated-User")
	}

	if userEmail == "" {
		return nil
	}
	user := User(userEmail)
	return &user
}

func TrivialAuth(db *ent.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user := ExtractUser(r.Header)

			if user == nil {
				http.Error(w, "401 Unauthorized: Authentication Required", http.StatusUnauthorized)
				return
			}

			_, err := db.Utente.Query().Where(utente.ID(string(*user))).Only(r.Context())

			if ent.IsNotFound(err) {
				http.Error(w, "401 Unauthorized: Authentication Required", http.StatusUnauthorized)
				return
			}

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Printf("Errore gestione utente: %v", err)
			}

			ctx := context.WithValue(r.Context(), UserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(r *http.Request) *User {
	if r == nil {
		log.Print("No request")
		return nil
	}
	user, ok := r.Context().Value(UserKey).(*User)
	if !ok {
		log.Print("Error getting user")
		return nil
	}

	return user
}
