package service

import (
	"database/sql"
	"fmt"
	"log"
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


func SaveGroupMessage(senderID int, groupID int, content string) error {
	log.Println("Memulai transaksi database untuk SaveGroupMessage")

	tx, err := config.DB.Begin()
	if err != nil {
		log.Println("Gagal memulai transaksi:", err)
		return err
	}
	log.Println("Transaksi dimulai")

	messageID, err := insertGroupMessage(tx, senderID, content, groupID)
	if err != nil {
		tx.Rollback()
		return err
	}

	receiverIDs, err := getGroupMemberIDs(tx, groupID, senderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	success := insertMessageReceivers(tx, messageID, receiverIDs)

	if err := tx.Commit(); err != nil {
		log.Println("Gagal commit transaksi:", err)
		return err
	}

	if !success {
		log.Println("Transaksi berhasil tapi beberapa receiver gagal disisipkan")
	}
	return nil
}

func insertGroupMessage(tx *sql.Tx, senderID int, content string, groupID int) (int, error) {
	log.Println("Menyisipkan pesan ke tabel group_messages")
	result, err := tx.Exec("INSERT INTO group_messages (sender_id, content, group_id) VALUES (?, ?, ?)", senderID, content, groupID)
	if err != nil {
		log.Println("Gagal menyisipkan pesan:", err)
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Println("Gagal mengambil LastInsertId:", err)
		return 0, err
	}
	log.Printf("ID pesan yang disisipkan: %d", lastID)
	return int(lastID), nil
}

func getGroupMemberIDs(tx *sql.Tx, groupID int, excludeSenderID int) ([]int, error) {
	log.Printf("Mengambil anggota group_id: %d", groupID)
	rows, err := tx.Query("SELECT user_id FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		log.Println("Gagal mengambil anggota grup:", err)
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Println("Gagal scan user_id:", err)
			continue
		}
		if id != excludeSenderID {
			ids = append(ids, id)
		}
	}
	if err := rows.Err(); err != nil {
		log.Println("Error saat iterasi rows:", err)
		return nil, err
	}
	return ids, nil
}

func insertMessageReceivers(tx *sql.Tx, messageID int, receiverIDs []int) bool {
	success := true
	for _, receiverID := range receiverIDs {
		log.Printf("Menyisipkan receiver_id: %d ke message_receivers", receiverID)
		_, err := tx.Exec("INSERT INTO message_receivers (message_id, receiver_id) VALUES (?, ?)", messageID, receiverID)
		if err != nil {
			log.Printf("Gagal menyisipkan receiver_id %d: %v", receiverID, err)
			success = false
		}
	}
	return success
}

func GetGroupMessages(group int) ([]dto.GetGroupMessage, error) {
	query := `
		SELECT id, sender_id, group_id, content, sent_at 
		FROM group_messages 
		WHERE group_id = ?
		ORDER BY sent_at ASC
	`

	rows, err := config.DB.Query(query, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []dto.GetGroupMessage
	for rows.Next() {
		var msg dto.GetGroupMessage
		var sentAtStr string

		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.GroupID, &msg.Content, &sentAtStr)
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