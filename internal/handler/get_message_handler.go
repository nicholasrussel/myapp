package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasrussel/myapp/internal/service"
)

func GetMessagesHandler(c *gin.Context) {
	roomIDStr := c.Query("chat_room_id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat_room_id parameter"})
		return
	}

	messages, err := service.GetMessagesByChatRoom(roomID)
	if err != nil {
		log.Println("DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

