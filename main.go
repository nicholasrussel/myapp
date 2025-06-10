package main

import (
	"fmt"
	"log"

	// Tambahkan ini
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/email"
	"github.com/nicholasrussel/myapp/internal/handler"
	"github.com/nicholasrussel/myapp/internal/service"
)

func main() {
	fmt.Println("Hello, Go!")

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load file .env")
	}

	email.InitGmailService()

	// Init DB
	config.InitDB()
	log.Println("✅ Koneksi ke database berhasil!")

	// Inisialisasi Gin router
	router := gin.Default()
	router.StaticFile("/", "index.html")

	// Routes
	log.Println("Mendaftarkan endpoint /login")

	router.POST("/login", service.RateLimitLoginMiddleware(), handler.LoginHandler)
	router.GET("/check-login", service.Authenticate(handler.CheckLoginHandler, 1))
	router.POST("/refresh", handler.RefreshTokenHandler)

	router.POST("/register", handler.RegisterHandler)
	router.GET("/logout", handler.Logout)

	router.GET("/2fa/setup/:username", handler.Setup2FAHandler)
	router.POST("/2fa/verify", handler.Verify2FAHandler)

	router.POST("/send-email", handler.SendEmailHandler)

	router.GET("/ws", handler.WebSocketHandler)

	router.POST("/messages", handler.SendMessageHandler)
	router.GET("/messages", handler.GetMessagesHandler)

	// Jalankan server
	svrPort := config.LoadEnv("SVR_PORT")
	log.Println("Connected to port " + svrPort)
	router.Run(":" + svrPort)
}
