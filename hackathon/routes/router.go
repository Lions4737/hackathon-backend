package routes

import (
	"github.com/gorilla/mux"
	"hackathon/handlers"
	"hackathon/middleware"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// session関連のルート
	r.HandleFunc("/api/sessionLogin", middleware.WithCORS(handlers.SessionLoginHandler)).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/sessionLogout", middleware.WithCORS(handlers.SessionLogoutHandler)).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/checkSession", middleware.WithCORS(middleware.RequireAuth(handlers.CheckSession))).Methods("OPTIONS", "GET")

	// posts関連のルート
	r.HandleFunc("/api/posts", middleware.WithCORS(middleware.RequireAuth(handlers.CreatePostHandler))).Methods("OPTIONS", "POST", "GET")
	r.HandleFunc("/api/all-posts", middleware.WithCORS(middleware.RequireAuth(handlers.GetAllPostsHandler))).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/my-posts", middleware.WithCORS(middleware.RequireAuth(handlers.GetMyPostsHandler))).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/posts/{id:[0-9]+}", middleware.WithCORS(middleware.RequireAuth(handlers.GetPostByID))).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/posts/{id:[0-9]+}/replies", middleware.WithCORS(middleware.RequireAuth(handlers.GetRepliesByPostID))).Methods("OPTIONS", "GET")
	r.HandleFunc("/api/posts/{id}/factcheck", middleware.WithCORS(handlers.FactCheckHandler)).Methods("OPTIONS", "GET")

	// db関連のルート
	r.HandleFunc("/api/registerUser", middleware.WithCORS(handlers.RegisterUserHandler)).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/me", middleware.WithCORS(middleware.RequireAuth(handlers.GetCurrentUser))).Methods("GET", "OPTIONS")

	// いいね関連のルート
	r.HandleFunc("/api/posts/{id:[0-9]+}/like", middleware.WithCORS(middleware.RequireAuth(handlers.LikePost))).Methods("OPTIONS", "POST")
	r.HandleFunc("/api/posts/{id:[0-9]+}/unlike", middleware.WithCORS(middleware.RequireAuth(handlers.UnlikePost))).Methods("OPTIONS", "DELETE")
	r.HandleFunc("/api/my-likes", middleware.WithCORS(middleware.RequireAuth(handlers.GetMyLikes))).Methods("OPTIONS", "GET")

	// プロフィール関連のルート
	r.HandleFunc("/api/profile", middleware.WithCORS(middleware.RequireAuth(handlers.GetProfileHandler))).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/profile", middleware.WithCORS(middleware.RequireAuth(handlers.UpdateProfileHandler))).Methods("POST", "OPTIONS")


	return r
}
