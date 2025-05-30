package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
	"github.com/nicholasrussel/myapp/internal/model"
	"github.com/nicholasrussel/myapp/internal/service"
)

func LoginHandler(c *gin.Context) {
	log.Println("LoginHandler dipanggil")

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Cek rate limit login gagal per email
	if service.IsLoginFailureLimited(req.Email, c) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Terlalu banyak login gagal. Coba lagi nanti."})
		return
	}

	// Coba login
	user, err := service.Login(req.Email, req.Password)
	if err != nil {
		// Jika gagal, catat sebagai login gagal
		service.MarkLoginFailure(req.Email, c)

		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Jika berhasil login, buat token
	token, err := service.GenerateToken(user.ID, user.Username, user.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Set cookie untuk Web
	c.SetCookie(service.TokenName, token, 60, "/", "", true, true)

	// Generate refresh token
	refreshToken, err := service.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat refresh token"})
		return
	}
	c.SetCookie(service.RefreshTokenName, refreshToken, 60, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"token":        token,
		"refreshToken": refreshToken,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token tidak ditemukan"})
		return
	}

	jwtKey := []byte(config.LoadEnv("JWT_REFRESH_KEY"))
	claims := &model.JWTClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token tidak valid"})
		return
	}

	newAccessToken, err := service.GenerateToken(claims.ID, claims.Username, claims.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal buat access token"})
		return
	}

	c.SetCookie(service.TokenName, newAccessToken, 3600 /* detik */, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"token": newAccessToken,
	})
}

func CheckLoginHandler(c *gin.Context) {
	c.String(http.StatusOK, "Selamat datang, user premium!")
}

func Logout(c *gin.Context) {
	service.ResetUserToken(c)
	log.Println("(SUCCESS)\t Logout request")
}

