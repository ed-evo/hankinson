package api_middlewares

import (
	"context"
	"log"
	"net/http"
)

type contextKey string

type User string

const UserKey contextKey = "userEmail"

func TrivialAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userEmail := r.Header.Get("Cf-Access-Authenticated-User-Email")

		if userEmail == "" {
			userEmail = r.Header.Get("X-Authenticated-User")
		}

		if userEmail == "" {
			http.Error(w, "401 Unauthorized: Authentication Required", http.StatusUnauthorized)
			return // Short-circuit: stop execution right here
		}

		user := User(userEmail)

		ctx := context.WithValue(r.Context(), UserKey, &user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
