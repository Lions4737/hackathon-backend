package main

import (
	"log"
	"net/http"
	"github.com/joho/godotenv"

	"hackathon/routes"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .env ファイルが読み込めませんでした")
	}
}

func main() {
	r := routes.SetupRouter()

	log.Println("Server listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
