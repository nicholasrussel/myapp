package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
	"github.com/nicholasrussel/myapp/internal/service"
)


func SendMessageHandler(c *gin.Context) {
	var msg dto.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var err error
	if msg.GroupID != 0 {
		log.Println("masuk if group id ada nilai:")
		err = service.SaveGroupMessage(msg.SenderID, msg.GroupID, msg.Content)
	} else {
		log.Println("masuk if chat personal:")
		err = service.SaveMessage(msg.SenderID, msg.ReceiverID, msg.Content)
	}

	if err != nil {
		log.Println("send messages error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}


