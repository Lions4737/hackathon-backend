package db

import (
	"database/sql"
	"fmt"
	"os"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	user := os.Getenv("MYSQL_USER")
	pwd := os.Getenv("MYSQL_PWD")
	socket := os.Getenv("MYSQL_HOST")
	dbname := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@%s/%s", user, pwd, socket, dbname)
	log.Println("接続先DSN:", dsn)
	return sql.Open("mysql", dsn)
}
