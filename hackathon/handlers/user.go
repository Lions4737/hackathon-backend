package handlers

import (
	"fmt"
	"net/http"

	"hackathon/db"
	"hackathon/model"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	if err := db.GetDB().Find(&users).Error; err != nil {
		http.Error(w, "ユーザー取得失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		fmt.Fprintf(w, "id: %s, username: %s, description: %s\n",
			user.ID, user.Username, user.Description)
	}
}
