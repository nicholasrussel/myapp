package service

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pquerna/otp/totp"
)

var (
	userSecrets = make(map[string]string)
	mu          sync.Mutex
	secretsFile = "secrets.json"
)

// Load secrets saat start
func init() {
	loadSecrets()
}

// Simpan secrets ke file
func saveSecrets() {
	mu.Lock()
	defer mu.Unlock()
	data, err := json.MarshalIndent(userSecrets, "", "  ")
	if err != nil {
		log.Println("Gagal menyimpan secrets:", err)
		return
	}
	os.WriteFile(secretsFile, data, 0644)
}

// Load secrets dari file
func loadSecrets() {
	data, err := os.ReadFile(secretsFile)
	if err != nil {
		log.Println("Belum ada file secrets, akan dibuat saat diperlukan")
		return
	}
	json.Unmarshal(data, &userSecrets)
}

// Generate secret dan URL OTP
func GenerateTOTPSecret(username string) (secret string, otpURL string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MyApp",
		AccountName: username,
	})
	if err != nil {
		return "", "", err
	}

	userSecrets[username] = key.Secret()
	saveSecrets()
	return key.Secret(), key.URL(), nil
}

// Validasi OTP
func ValidateTOTPCode(username, code string) bool {
	secret, ok := userSecrets[username]
	if !ok {
		log.Printf("Secret tidak ditemukan untuk user: %s", username)
		return false
	}

	now := time.Now()
	valid := totp.Validate(code, secret)
	if !valid {
		log.Printf("OTP tidak valid untuk user: %s, code: %s, secret: %s, time: %s", username, code, secret, now.Format(time.RFC3339))
	}
	return valid
}
