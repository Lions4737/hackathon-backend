package middleware

import (
	"context"
	"net/http"

	"hackathon/firebase"
)

type contextKey string

const UserIDKey = contextKey("userID")

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		client, _ := firebase.GetAuthClient()
		token, err := client.VerifyIDToken(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, token.UID)
		next(w, r.WithContext(ctx))
	}
}
