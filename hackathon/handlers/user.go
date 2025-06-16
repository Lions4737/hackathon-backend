package handlers

import (
	"fmt"
	"net/http"
	"hackathon/db"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := db.Connect()
	if err != nil {
		http.Error(w, "DB接続エラー: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT id, name, age FROM user")
	if err != nil {
		http.Error(w, "クエリエラー: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			http.Error(w, "スキャンエラー: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "id: %s, name: %s, age: %d\n", id, name, age)
	}
}
