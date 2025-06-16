package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"hackathon/firebase"
)

func SessionLoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	idToken := r.FormValue("idToken")
	if idToken == "" {
		http.Error(w, "ID token is required", http.StatusBadRequest)
		return
	}

	app, err := firebase.InitFirebaseApp()
	if err != nil {
		http.Error(w, "Firebase init error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	auth, err := app.Auth(context.Background())
	if err != nil {
		http.Error(w, "Auth error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = auth.VerifyIDToken(context.Background(), idToken)
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
