package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" // これを忘れずに！
	"hackathon/db"
	"hackathon/middleware"
	"hackathon/model"
)

// GET /api/posts
func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	database := db.GetDB()

	var user model.User
	if err := database.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var posts []model.Post
	if err := database.Preload("User").
		Where("user_id != ?", user.ID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GET /api/my-posts
func GetMyPostsHandler(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	database := db.GetDB()

	var user model.User
	if err := database.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var posts []model.Post
	if err := database.Preload("User").
		Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GET /api/users/{id}/posts
func GetPostsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	database := db.GetDB()

	var posts []model.Post
	if err := database.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetUserProfileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	database := db.GetDB()

	var user model.User
	if err := database.First(&user, userID).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	profile := map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"description":   user.Description,
		"profile_image": user.ProfileImage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}