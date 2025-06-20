package handlers

import (
  "encoding/json"
  "net/http"
  "log"

  "hackathon/db"
  "hackathon/middleware"
  "hackathon/model"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
  uidVal := r.Context().Value(middleware.UserIDKey)
  uid, ok := uidVal.(string)
  if !ok || uid == "" {
    http.Error(w, "unauthorized", http.StatusUnauthorized)
    return
  }

  var user model.User
  if err := db.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
    log.Println("⚠️ ユーザー取得失敗:", uid)
    http.Error(w, "user not found", http.StatusNotFound)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(user)
}
