package service

import (
	"github.com/nicholasrussel/myapp/config"
)

func CreateGroup(name string, createdBy int, memberIDs []int) (int, error) {
	tx, err := config.DB.Begin()
	if err != nil {
		return 0, err
	}

	// Step 1: Insert group
	res, err := tx.Exec("INSERT INTO `groups` (name, created_by) VALUES (?, ?)", name, createdBy)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	groupID64, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	groupID := int(groupID64)

	// Step 2: Insert members
	stmt, err := tx.Prepare("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	for _, userID := range memberIDs {
		_, err := stmt.Exec(groupID, userID)
		if err != nil {
			tx.Rollback()
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
