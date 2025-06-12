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

type Client struct {
	conn   *websocket.Conn
	userID string
}

var (
	clients  = make(map[string]*Client)
	lock     sync.Mutex
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func WebSocketHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
		return
	}

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
	lock.Lock()
	clients[userID] = client
	lock.Unlock()

	log.Printf("User %s connected", userID)

	for {
		var msg dto.Message
		if err := ws.ReadJSON(&msg); err != nil {
			log.Printf("Read error from user %s: %v", userID, err)
			break
		}
		handleMessage(msg)
	}
}

func handleMessage(msg dto.Message) {
	if msg.GroupID > 0 {
		handleGroupMessage(msg)
	} else {
		handlePersonalMessage(msg)
	}
}

func handleGroupMessage(msg dto.Message) {
	err := service.SaveGroupMessage(msg.SenderID, msg.GroupID, msg.Content)
	if err != nil {
		log.Printf("Failed to save group message: %v", err)
		return
	}

	memberIDs, err := service.GetGroupMembers(msg.GroupID)
	if err != nil {
		log.Printf("Failed to get group members: %v", err)
		return
	}

	for _, uid := range memberIDs {
		if uid != msg.SenderID {
			sendToClient(uid, msg)
		}
	}
}

func handlePersonalMessage(msg dto.Message) {
	err := service.SaveMessage(msg.SenderID, msg.ReceiverID, msg.Content)
	if err != nil {
		log.Printf("Failed to save personal message: %v", err)
		return
	}
	sendToClient(msg.ReceiverID, msg)
}

func sendToClient(userID int, msg dto.Message) {
	uidStr := strconv.Itoa(userID)
	lock.Lock()
	receiverClient, ok := clients[uidStr]
	lock.Unlock()

	if ok {
		err := receiverClient.conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Write error to user %d: %v", userID, err)
		}
	}
}
