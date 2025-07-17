package dto

import "time"

type Message struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID *int    `json:"receiver_id"`
	GroupID    *int    `json:"group_id"`
	Content    string `json:"content"`
	ChatRoomID *int 	  `json:"chat_room_id"`
}

type GetMessage struct {
	ID		   int	  		`json:"id"`
	SenderID   int    		`json:"sender_id"`
	ReceiverID int    		`json:"receiver_id"`
	Content    string 		`json:"content"`
	SentAt	   time.Time 	`json:"sent_at"`
}