package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"

	"hackathon/db"
	"hackathon/middleware"
	"hackathon/model"
)

type CreatePostRequest struct {
	Content      string `json:"content"`
	ParentPostID *uint  `json:"parent_post_id,omitempty"` // 🔥 リプライ対応
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
		UserID:       user.ID,
		Content:      req.Content,
		IsReply:      req.ParentPostID != nil,
		ParentPostID: req.ParentPostID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 💡 トランザクションで作成＋カウント更新
	err := db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&post).Error; err != nil {
			return err
		}
		// 🔁 親ポストがある場合、reply_count を +1
		if req.ParentPostID != nil {
			if err := tx.Model(&model.Post{}).
				Where("id = ?", *req.ParentPostID).
				Update("reply_count", gorm.Expr("reply_count + 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
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
