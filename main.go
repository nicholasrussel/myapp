package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/handler"
	"github.com/nicholasrussel/myapp/internal/service"
)

func main() {
	fmt.Println("Hello, Go!")


	err := godotenv.Load()
	if err != nil {
		log.Fatal("Gagal load file .env")
	}

	// Init DB
	config.InitDB()
	log.Println("✅ Koneksi ke database berhasil!")
	log.Println("Mendaftarkan endpoint /login")
	
	router := mux.NewRouter()

	router.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	router.HandleFunc("/check-login", service.Authenticate(handler.CheckLoginHandler, 1)).Methods("GET")
	router.HandleFunc("/refresh", handler.RefreshTokenHandler).Methods("POST")

	
	router.HandleFunc("/logout", handler.Logout).Methods("GET")

	log.Println("Endpoint /login terdaftar")
	svrPort := config.LoadEnv("SVR_PORT")
	log.Println("Connected to port " + svrPort)
	addr := ":" + svrPort
	http.ListenAndServe(addr, router)
}

