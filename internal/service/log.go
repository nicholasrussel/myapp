package service

import (
	"log"

	"github.com/nicholasrussel/myapp/config"
)

func InsertLog(userID *int, actionType, status, message, ip, userAgent string) {
	query := `INSERT INTO logs (user_id, action_type, status, message, ip_address, user_agent) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := config.DB.Exec(query, userID, actionType, status, message, ip, userAgent)
	if err != nil {
		log.Println("Gagal insert log:", err)
	}
}
