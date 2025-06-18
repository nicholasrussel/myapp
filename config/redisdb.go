package config

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "", // default tanpa password
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Gagal konek ke Redis: %v", err)
	}

	log.Println("✅ Koneksi ke Redis berhasil!")
}
