package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	// 環境変数から DB 情報取得
	user := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PWD")
	dbname := os.Getenv("MYSQL_DATABASE")
	connName := os.Getenv("MYSQL_CONNECTION_NAME") // 追加（例: "project:region:instance"）

	// DSNの組み立て（Unixソケット経由）
	dsn := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", user, pwd, connName, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, "DB接続エラー: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, age FROM user")
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

func main() {
	http.HandleFunc("/users", userHandler)
	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
