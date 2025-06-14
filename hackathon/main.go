package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// CORS対応のラッパー
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 必須ヘッダーを設定
		w.Header().Set("Access-Control-Allow-Origin", "*") // 本番は "*" を制限すべき
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// プリフライトリクエスト（OPTIONS）の場合はすぐ返す
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 通常のリクエスト処理
		next(w, r)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	user := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PWD")
	socket := os.Getenv("MYSQL_HOST") // unix(/cloudsql/...) 形式
	dbname := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@%s/%s", user, pwd, socket, dbname)
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
	http.HandleFunc("/users", withCORS(userHandler))
	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
