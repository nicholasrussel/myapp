package service

import (
	"github.com/nicholasrussel/myapp/config"
)

func CreateGroup(name string, createdBy int, memberIDs []int) (int, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.Exec("INSERT INTO `groups` (name, created_by) VALUES (?, ?)", name, createdBy)
	if err != nil {
		return 0, err
	}
	groupID64, _ := res.LastInsertId()
	groupID := int(groupID64)

	res, err = tx.Exec("INSERT INTO chat_rooms (is_group, name) VALUES (?, ?)", true, name)
	if err != nil {
		return 0, err
	}
	chatRoomID64, _ := res.LastInsertId()
	chatRoomID := int(chatRoomID64)

	stmtChatRoom, err := tx.Prepare("INSERT INTO chat_room_users (chat_room_id, user_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmtChatRoom.Close()

	allMembers := append(memberIDs, createdBy)
	for _, uid := range allMembers {
		if _, err := stmtChatRoom.Exec(chatRoomID, uid); err != nil {
			return 0, err
		}
	}

	stmtGroup, err := tx.Prepare("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmtGroup.Close()

	for _, uid := range allMembers {
		if _, err := stmtGroup.Exec(groupID, uid); err != nil {
			return 0, err
		}
	}

	return groupID, tx.Commit()
}


func GetGroupMembers(groupID int) ([]int, error) {
	rows, err := config.DB.Query("SELECT user_id FROM group_members WHERE group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}
	return members, nil
}
