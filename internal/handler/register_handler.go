package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
	"github.com/nicholasrussel/myapp/internal/service"
)

func RegisterHandler(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		service.InsertLog(nil, "REGISTER", "FAILED", "Invalid request body", c.ClientIP(), c.Request.UserAgent())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	err := service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		service.InsertLog(nil, "REGISTER", "FAILED", err.Error(), c.ClientIP(), c.Request.UserAgent())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.Login(req.Email, req.Password)
	if err != nil {
		service.InsertLog(nil, "REGISTER_LOGIN", "FAILED", "Gagal login setelah registrasi", c.ClientIP(), c.Request.UserAgent())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal login setelah registrasi"})
		return
	}

	service.InsertLog(&user.ID, "REGISTER", "SUCCESS", "Registrasi dan login berhasil", c.ClientIP(), c.Request.UserAgent())

	CreateTokenHandler(c, user)
}


