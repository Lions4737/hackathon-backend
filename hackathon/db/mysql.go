package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PWD")
	host := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")

	//dsn := fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FTokyo", user, pass, host, dbname)
	dsn := fmt.Sprintf("%s:%s@%s/%s", user, pwd, socket, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("GORM DB接続失敗: %w", err)
	}

	log.Println("✅ GORM: DB接続成功")
	DB = db
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
