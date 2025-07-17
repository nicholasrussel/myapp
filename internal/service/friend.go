package service

import (
	"errors"

	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/constants"
)

func SendFriendRequest(fromID, toID int) error {
	if fromID == toID {
		return errors.New("cannot send friend request to yourself")
	}

	var exists int
	err := config.DB.QueryRow(`
		SELECT 1 FROM friend_requests 
		WHERE from_user_id = ? AND to_user_id = ?
	`, fromID, toID).Scan(&exists)

	if err == nil {
		return errors.New("friend request already sent")
	}

	_, err = config.DB.Exec(`
		INSERT INTO friend_requests (from_user_id, to_user_id)
		VALUES (?, ?)
	`, fromID, toID)

	return err
}

func GetFriendRequests(userID int) ([]map[string]interface{}, error) {
	rows, err := config.DB.Query(`
		SELECT id, from_user_id, status, created_at 
		FROM friend_requests 
		WHERE to_user_id = ? AND status = 'pending'
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, fromUser int
		var status string
		var createdAt string

		err := rows.Scan(&id, &fromUser, &status, &createdAt)
		if err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":            id,
			"from_user_id":  fromUser,
			"status":        status,
			"created_at":    createdAt,
		})
	}

	return results, nil
}

func HandleFriendRequestAction(userID, fromUserID int, action string) error {
	if action == constants.FriendRequestAccept {

		_, err := config.DB.Exec(`
			UPDATE friend_requests 
			SET status = 'accepted'
			WHERE from_user_id = ? AND to_user_id = ?
		`, fromUserID, userID)
		if err != nil {
			return err
		}

		user1 := min(fromUserID, userID)
		user2 := max(fromUserID, userID)

		_, err = config.DB.Exec(`
			INSERT INTO friends (user1_id, user2_id) VALUES (?, ?)
		`, user1, user2)
		return err

	} else if action == constants.FriendRequestReject {

		_, err := config.DB.Exec(`
			UPDATE friend_requests 
			SET status = 'rejected'
			WHERE from_user_id = ? AND to_user_id = ?
		`, fromUserID, userID)
		return err

	}

	return errors.New("invalid action")
}

func BlockFriend(userID, friendID int) error {
	var exists int
	err := config.DB.QueryRow(`
		SELECT 1 FROM friends 
		WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)
	`, userID, friendID, friendID, userID).Scan(&exists)

	if err != nil {
		return errors.New("you are not friends")
	}

	_, err = config.DB.Exec(`
		UPDATE friends 
		SET is_blocked_by = ?
		WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)
	`, userID, userID, friendID, friendID, userID)

	return err
}

func DeleteFriend(userID, friendID int) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM friends 
		WHERE (user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)
	`, userID, friendID, friendID, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM friend_requests
		WHERE ((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?))
		  AND status = 'accepted'
	`, userID, friendID, friendID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}


