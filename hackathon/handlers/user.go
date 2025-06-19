// handlers/user.go

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"hackathon/db"
	"hackathon/firebase"
	"hackathon/model"
)

type RegisterUserRequest struct {
	IDToken string `json:"idToken"`
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "リクエスト形式が不正です", http.StatusBadRequest)
		return
	}

	authClient, err := firebase.GetAuthClient()
	if err != nil {
		http.Error(w, "Firebase init error", http.StatusInternalServerError)
		return
	}

	// Verify ID token
	token, err := authClient.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		http.Error(w, "IDトークンが無効です", http.StatusUnauthorized)
		return
	}

	uid := token.UID
	dbConn := db.GetDB()

	// すでに登録されていれば無視
	var existing model.User
	if err := dbConn.Where("firebase_uid = ?", uid).First(&existing).Error; err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("既にユーザーが存在します"))
		return
	}

	user := model.User{
		FirebaseUID:  uid,
		Username:     "Guest",          // 初期値
		Description:  "",         // 初期値
		ProfileImage: "default.png",        // 初期値
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := dbConn.Create(&user).Error; err != nil {
		http.Error(w, "ユーザー作成に失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ユーザー登録完了"))
}
