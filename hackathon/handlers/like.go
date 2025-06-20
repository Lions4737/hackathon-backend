package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"hackathon/db"
	"hackathon/middleware"
	"hackathon/model"
	"gorm.io/gorm"
)

// POST /api/posts/{id}/like
func LikePost(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := mux.Vars(r)["id"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := db.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var like model.Like
	if err := db.DB.Where("user_id = ? AND post_id = ?", user.ID, postID).First(&like).Error; err == nil {
		http.Error(w, "already liked", http.StatusConflict)
		return
	}

	like = model.Like{UserID: user.ID, PostID: uint(postID)}
	if err := db.DB.Create(&like).Error; err != nil {
		http.Error(w, "failed to like", http.StatusInternalServerError)
		return
	}

	db.DB.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1))

	w.WriteHeader(http.StatusCreated)
}

// DELETE /api/posts/{id}/unlike
func UnlikePost(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postIDStr := mux.Vars(r)["id"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	var user model.User
	if err := db.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if err := db.DB.Where("user_id = ? AND post_id = ?", user.ID, postID).Delete(&model.Like{}).Error; err != nil {
		http.Error(w, "failed to unlike", http.StatusInternalServerError)
		return
	}

	db.DB.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1))

	w.WriteHeader(http.StatusOK)
}

// GET /api/my-likes
func GetMyLikes(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user model.User
	if err := db.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var likes []model.Like
	db.DB.Where("user_id = ?", user.ID).Find(&likes)

	postIDs := make([]uint, len(likes))
	for i, like := range likes {
		postIDs[i] = like.PostID
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postIDs)
}
