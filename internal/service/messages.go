package service

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/handler/dto"
)

func SaveMessage(senderID, chatRoomID int, content string) error {
	query := `
		INSERT INTO messages (chat_room_id, sender_id, content, sent_at)
		VALUES (?, ?, ?, ?)
	`
	_, err := config.DB.Exec(query, chatRoomID, senderID, content, time.Now())
	return err
}

func FindOrCreateChatRoom(user1, user2 int) (int, error) {
	var chatRoomID int

	// Pastikan user1 selalu lebih kecil dari user2 untuk konsistensi (opsional)
	if user1 > user2 {
		user1, user2 = user2, user1
	}

	// Cek apakah chat room 1-on-1 sudah ada (regardless of user order)
	query := `
	SELECT cru.chat_room_id
	FROM chat_room_users cru
	JOIN chat_rooms cr ON cr.id = cru.chat_room_id
	WHERE cr.is_group = FALSE AND cru.chat_room_id IN (
		SELECT cru2.chat_room_id
		FROM chat_room_users cru2
		WHERE cru2.user_id IN (?, ?)
		GROUP BY cru2.chat_room_id
		HAVING COUNT(DISTINCT cru2.user_id) = 2
	)
	GROUP BY cru.chat_room_id
	LIMIT 1
	`

	err := config.DB.QueryRow(query, user1, user2).Scan(&chatRoomID)
	if err == nil {
		return chatRoomID, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	// Belum ada, buat chat room baru
	tx, err := config.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	insertRoomQuery := `INSERT INTO chat_rooms (is_group) VALUES (FALSE)`
	res, err := tx.Exec(insertRoomQuery)
	if err != nil {
		return 0, err
	}

	roomID64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	chatRoomID = int(roomID64)

	// Tambahkan kedua user ke chat_room_users
	insertUserQuery := `INSERT INTO chat_room_users (chat_room_id, user_id) VALUES (?, ?), (?, ?)`
	_, err = tx.Exec(insertUserQuery, chatRoomID, user1, chatRoomID, user2)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return chatRoomID, nil
}


func GetMessagesByChatRoom(chatRoomID int) ([]dto.GetMessage, error) {
	query := `
		SELECT id, sender_id, content, sent_at 
		FROM messages 
		WHERE chat_room_id = ?
		ORDER BY sent_at ASC
	`

	rows, err := config.DB.Query(query, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []dto.GetMessage
	for rows.Next() {
		var msg dto.GetMessage
		var sentAtStr string

		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.Content, &sentAtStr)
		if err != nil {
			return nil, err
		}

		msg.SentAt, err = time.Parse("2006-01-02 15:04:05", sentAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sent_at: %v", err)
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

func GetGroupMemberIDs(tx *sql.Tx, groupID int, excludeSenderID int) ([]int, error) {
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

func InsertMessageReceivers(tx *sql.Tx, messageID int, receiverIDs []int) error {
	for _, receiverID := range receiverIDs {
		log.Printf("Menyisipkan receiver_id: %d ke message_receivers", receiverID)
		_, err := tx.Exec("INSERT INTO message_receivers (message_id, receiver_id) VALUES (?, ?)", messageID, receiverID)
		if err != nil {
			log.Printf("Gagal menyisipkan receiver_id %d: %v", receiverID, err)
			return nil
		}
	}
	return nil
}

func GetChatRoomMemberIDs(chatRoomID int) ([]int, error) {
	rows, err := config.DB.Query("SELECT user_id FROM chat_room_users WHERE chat_room_id = ?", chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
