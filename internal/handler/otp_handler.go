package handler

import (
	"github.com/nicholasrussel/myapp/internal/service"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup2FAHandler(c *gin.Context) {
	username := c.Param("username")

	secret, otpURL, err := service.GenerateTOTPSecret(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal generate secret"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret":       secret,
		"otp_auth_url": otpURL,
	})
}

func Verify2FAHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	if service.ValidateTOTPCode(req.Username, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "Berhasil verifikasi 2FA"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kode OTP tidak valid"})
	}
}
