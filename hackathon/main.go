package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"hackathon/db"
	"hackathon/model"
	"hackathon/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env ファイルが読み込めませんでした")
	}
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 マイグレーション実行中...")
	if err := db.GetDB().AutoMigrate(&model.User{}, &model.Post{}, &model.Like{}); err != nil {
		log.Println("⚠️ ユーザー作成失敗（重複の可能性あり）:", err)
	}
	log.Println("✅ マイグレーション完了")

	r := routes.SetupRouter()

	// ✅ Cloud Run 互換：PORT 環境変数を利用
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル開発用のデフォルト
	}
	log.Printf("🚀 Server listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
