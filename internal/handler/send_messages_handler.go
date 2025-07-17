package handler

import (
	"log"
	"net/http"
	"strconv"

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

	var chatRoomID int
	var err error

	if msg.ChatRoomID != nil && *msg.ChatRoomID > 0 {
		chatRoomID = *msg.ChatRoomID
	} else if msg.ReceiverID != nil && *msg.ReceiverID > 0 {
		// Buat atau cari chat room berdasarkan sender & receiver
		chatRoomID, err = service.FindOrCreateChatRoom(msg.SenderID, *msg.ReceiverID)
		if err != nil {
			log.Println("error finding/creating chat room:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat room"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either chat_room_id or receiver_id is required"})
		return
	}

	err = service.SaveMessage(msg.SenderID, chatRoomID, msg.Content)
	if err != nil {
		log.Println("send message error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Message sent successfully",
		"chat_room_id": chatRoomID,
	})
}

func CreateOrFindChatRoomHandler(c *gin.Context) {
	user1Str := c.Query("user1")
	user2Str := c.Query("user2")

	user1, err := strconv.Atoi(user1Str)
	if err != nil || user1 <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user1"})
		return
	}
	user2, err := strconv.Atoi(user2Str)
	if err != nil || user2 <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user2"})
		return
	}

	if user1 == user2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create chat room with self"})
		return
	}

	roomID, err := service.FindOrCreateChatRoom(user1, user2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create or find chat room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"chat_room_id": roomID,
	})
}

