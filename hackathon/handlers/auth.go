package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"encoding/json"

	"hackathon/firebase"
	"hackathon/middleware"
)

func SessionLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
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

	// Verify ID token
	_, err = authClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    idToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Login successful")
}


func SessionLogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Logout successful")
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(middleware.UserIDKey).(string)
	json.NewEncoder(w).Encode(map[string]string{"uid": uid})
}
