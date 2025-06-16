package routes

import (
	"github.com/gorilla/mux"
	"hackathon/handlers"
	"hackathon/middleware"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// APIルートの設定
	r.HandleFunc("/api/users", middleware.WithCORS(handlers.UserHandler)).Methods("GET")
	r.HandleFunc("/api/sessionLogin", middleware.WithCORS(handlers.SessionLoginHandler)).Methods("POST")
	r.HandleFunc("/api/sessionLogout", middleware.WithCORS(handlers.SessionLogoutHandler)).Methods("POST")

	return r
}
