package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicholasrussel/myapp/internal/service"
)

func GetMessagesHandler(c *gin.Context) {
	user1Str := c.Query("user1")
	user2Str := c.Query("user2")

	user1, err := strconv.Atoi(user1Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user1 parameter"})
		return
	}
	user2, err := strconv.Atoi(user2Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user2 parameter"})
		return
	}

	messages, err := service.GetMessagesBetweenUsers(user1, user2)
	if err != nil {
		log.Println("DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func GetGroupMessagesHandler(c *gin.Context) {
	groupStr := c.Query("group")

	group, err := strconv.Atoi(groupStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user1 parameter"})
		return
	}


	messages, err := service.GetGroupMessages(group)
	if err != nil {
		log.Println("DB error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
