package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"hackathon/firebase"
	"hackathon/middleware"
)

func isProduction() bool {
	return os.Getenv("ENV") == "production"
}

func newSessionCookie(value string, expiry time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    value,
		Path:     "/",
		Expires:  expiry,
		HttpOnly: true,
		Secure:   isProduction(), // 本番では https 用に true
	}
	if isProduction() {
		cookie.SameSite = http.SameSiteNoneMode // クロスドメインには必須
	} else {
		cookie.SameSite = http.SameSiteLaxMode // 開発環境では Lax が適切
	}
	return cookie
}

func SessionLoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	idToken := r.FormValue("idToken")
	if idToken == "" {
		http.Error(w, "ID token is required", http.StatusBadRequest)
		return
	}

	authClient, err := firebase.GetAuthClient()
	if err != nil {
		http.Error(w, "Firebase init error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify the ID token
	_, err = authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Set session cookie
	http.SetCookie(w, newSessionCookie(idToken, time.Now().Add(24*time.Hour)))

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Login successful")
}

func SessionLogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Expire the session cookie
	http.SetCookie(w, newSessionCookie("", time.Now().Add(-1*time.Hour)))

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Logout successful")
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(string)
	json.NewEncoder(w).Encode(map[string]string{"uid": uid})
}
