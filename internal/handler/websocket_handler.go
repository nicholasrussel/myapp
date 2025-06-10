package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
	"github.com/nicholasrussel/myapp/internal/service"
)

// Client represents a connected WebSocket client
type Client struct {
	conn   *websocket.Conn
	userID string
}

var (
	clients = make(map[string]*Client) // map userID to Client
	lock    sync.Mutex
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins
	}
)


func WebSocketHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer func() {
		ws.Close()
		log.Printf("User %s disconnected", userID)
		lock.Lock()
		delete(clients, userID)
		lock.Unlock()
	}()

	client := &Client{conn: ws, userID: userID}

	// Register client
	lock.Lock()
	clients[userID] = client
	lock.Unlock()

	log.Printf("User %s connected", userID)

	// Read messages loop
	for {
		var msg dto.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error from user %s: %v", userID, err)
			break
		}

		log.Printf("Received message from %s to %s: %s", msg.SenderID, msg.ReceiverID, msg.Content)

		err = service.SaveMessage(msg.SenderID, msg.ReceiverID, msg.Content)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
		}



		// Send message to receiver if online
		lock.Lock()
		receiverIDStr := strconv.Itoa(msg.ReceiverID)
		receiverClient, ok := clients[receiverIDStr]
		lock.Unlock()

		if ok {
			err = receiverClient.conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Write error to user %s: %v", msg.ReceiverID, err)
			}
		} else {
			log.Printf("User %s is offline, cannot deliver message", msg.ReceiverID)
		}
	}
}
