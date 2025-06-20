package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"hackathon/db"
	"hackathon/middleware"
	"hackathon/model"
)

type CreatePostRequest struct {
	Content string `json:"content"`
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("✅ context raw value: %v", r.Context().Value(middleware.UserIDKey))
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	uid := r.Context().Value(middleware.UserIDKey).(string)
	log.Printf("📛 UID from session: %s", uid)
	var user model.User
	if err := db.GetDB().Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	post := model.Post{
		UserID:    user.ID,
		Content:   req.Content,
		IsReply:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.GetDB().Create(&post).Error; err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// 作成した投稿に対応する User を preload で取得
	if err := db.GetDB().Preload("User").First(&post, post.ID).Error; err != nil {
		http.Error(w, "Failed to load post with user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}
