package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password string) (*model.User, error) {
	var user model.User
	err := config.DB.QueryRow("SELECT id, username, email, password, user_type FROM users WHERE email = ?", email).
			Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.UserType)

	if err == sql.ErrNoRows {
		return nil, errors.New("Email not found")
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("Password salah")
	}

	return &user, nil
}

func Register(username, email, password string) error {
	// Cek duplikat email
	var exists int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("Gagal memeriksa email")
	}
	if exists > 0 {
		return fmt.Errorf("Email sudah terdaftar")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Gagal meng-hash password")
	}

	// Simpan user
	_, err = config.DB.Exec(`INSERT INTO users (username, email, password, user_type) VALUES (?, ?, ?, 0)`,
		username, email, hashed)
	if err != nil {
		return fmt.Errorf("Gagal menyimpan user")
	}

	return nil
}
