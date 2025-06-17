package main

import (
	"log"
	"net/http"
	"os"

	"hackathon/db"
	"hackathon/model"
	"hackathon/routes"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal(err)
	}

	log.Println("🚀 マイグレーション実行中...")
	if err := db.GetDB().AutoMigrate(&model.User{}, &model.Post{}, &model.Like{}); err != nil {
		log.Fatalf("❌ マイグレーション失敗: %v", err)
	}
	log.Println("✅ マイグレーション完了")

	// 任意: 存在しない場合のみユーザー作成
	user := model.User{
		Username:    "taro",
		FirebaseUID: "uid123",
		Description: "テストユーザー",
	}
	db.GetDB().FirstOrCreate(&user, model.User{FirebaseUID: "uid123"})

	// Cloud Run対応: 環境変数からPORT取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := routes.SetupRouter()
	log.Printf("🚀 Server listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
