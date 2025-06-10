package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasrussel/myapp/internal/email"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
)

func SendEmailHandler(c *gin.Context) {
	var req dto.EmailRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := email.SendEmail(req.To, req.Subject, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
