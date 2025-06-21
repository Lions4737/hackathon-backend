package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"hackathon/db"
	"hackathon/model"
	"hackathon/middleware"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var found model.User
	if err := db.DB.Where("firebase_uid = ?", uid).First(&found).Error; err != nil {
		log.Println("❌ ユーザー取得失敗:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(found)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	uid, ok := uidVal.(string)
	if !ok || uid == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Username     string `json:"username"`
		Description  string `json:"description"`
		ProfileImage string `json:"profile_image"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := db.DB.Model(&model.User{}).Where("firebase_uid = ?", uid).
		Updates(model.User{
			Username:     input.Username,
			Description:  input.Description,
			ProfileImage: input.ProfileImage,
		}).Error; err != nil {
		log.Println("❌ プロフィール更新失敗:", err)
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
