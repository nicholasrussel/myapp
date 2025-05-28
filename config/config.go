package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
    dbHost := LoadEnv("DB_HOST")
    dbPort := LoadEnv("DB_PORT")
    dbUser := LoadEnv("DB_USER")
    dbPass := LoadEnv("DB_PASS")
    dbName := LoadEnv("DB_NAME")

	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASS:", dbPass)
	fmt.Println("DB_NAME:", dbName)

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

    var err error
    for i := 0; i < 10; i++ {
        DB, err = sql.Open("mysql", dsn)
        if err == nil {
            err = DB.Ping()
            if err == nil {
                fmt.Println("Berhasil konek ke DB!")
                return
            }
        }
        fmt.Printf("Gagal konek ke DB, coba lagi dalam 3 detik... (%d/10)\n", i+1)
        time.Sleep(3 * time.Second)
    }
    log.Fatal("Gagal konek ke DB setelah 10 kali percobaan:", err)
}



