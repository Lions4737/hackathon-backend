package routes

import (
	"github.com/gorilla/mux"
	"hackathon/handlers"
	"hackathon/middleware"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// APIルートの設定
	r.HandleFunc("/api/users", middleware.WithCORS(handlers.UserHandler)).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/sessionLogin", middleware.WithCORS(handlers.SessionLoginHandler)).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/sessionLogout", middleware.WithCORS(handlers.SessionLogoutHandler)).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/checkSession", middleware.WithCORS(middleware.RequireAuth(handlers.CheckSession))).Methods("OPTIONS", "GET")



	return r
}
