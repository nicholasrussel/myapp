package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func LoadEnv(key string) string {
	_ = godotenv.Load()
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("ENV %s tidak ditemukan", key)
	}
	return value
}

func InitDB() {
	dbPort := LoadEnv("DB_PORT")
	dbName := LoadEnv("DB_NAME")
	dsn := "root:@tcp(" + dbPort + ")/" + dbName

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Gagal buka koneksi DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal konek ke DB:", err)
	}
}
