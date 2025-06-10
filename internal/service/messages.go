package service

import (
	"fmt"
	"time"

	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
)

func SaveMessage(senderID, receiverID int, content string) error {
	query := "INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)"
	_, err := config.DB.Exec(query, senderID, receiverID, content)
	return err
}

func GetMessagesBetweenUsers(user1, user2 int) ([]dto.GetMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, content, sent_at 
		FROM messages 
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY sent_at ASC
	`

	rows, err := config.DB.Query(query, user1, user2, user2, user1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []dto.GetMessage
	for rows.Next() {
		var msg dto.GetMessage
		var sentAtStr string

		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &sentAtStr)
		if err != nil {
			return nil, err
		}

		// Parse string ke time.Time
		msg.SentAt, err = time.Parse("2006-01-02 15:04:05", sentAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sent_at: %v", err)
		}

		messages = append(messages, msg)
	}

	return messages, nil
}
