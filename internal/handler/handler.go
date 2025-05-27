package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
	"github.com/nicholasrussel/myapp/internal/model"
	"github.com/nicholasrussel/myapp/internal/service"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler dipanggil")
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := service.GenerateToken(user.ID, user.Username, user.UserType)
	if err != nil {
		http.Error(w, "Gagal membuat token", http.StatusInternalServerError)
		return
	}

	// Set cookie khusus untuk Web
	http.SetCookie(w, &http.Cookie{
		Name:     service.TokenName,
		Value:    token,
		Expires:  time.Now().Add(1 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	})

	// Buat refresh token
	refreshToken, err := service.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Gagal membuat refresh token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(1 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	})
	

	// Return response untuk mobile juga
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":  user,
		"token": token,
		"refreshToken": refreshToken,
	})
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "Refresh token tidak ditemukan", http.StatusUnauthorized)
		return
	}

	tokenStr := cookie.Value
	jwtKey := []byte(config.LoadEnv("JWT_REFRESH_KEY"))
	claims := &model.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Refresh token tidak valid", http.StatusUnauthorized)
		return
	}

	// Buat access token baru
	newAccessToken, err := service.GenerateToken(claims.ID, claims.Username, claims.UserType)
	if err != nil {
		http.Error(w, "Gagal buat access token", http.StatusInternalServerError)
		return
	}

	// Set cookie baru (jika dari browser)
	http.SetCookie(w, &http.Cookie{
		Name:     service.TokenName,
		Value:    newAccessToken,
		Expires:  time.Now().Add(60 * time.Minute),
		Secure:   true,
		HttpOnly: true,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": newAccessToken,
	})
}

func CheckLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Selamat datang, user premium!"))
}

func Logout(w http.ResponseWriter, r *http.Request) {

	service.ResetUserToken(w)
	log.Println("(SUCCESS)\t", "Logout request")
}
