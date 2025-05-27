package service

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/model"
)

var TokenName = "loginToken"

func GenerateToken(id int, username string, userType int) (string, error) {
	jwtKey := []byte(config.LoadEnv("JWT_KEY"))
	tokenExpiryTime := time.Now().Add(1 * time.Minute)

	claims := &model.JWTClaims{
		ID:       id,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenExpiryTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateRefreshToken(id int) (string, error) {
	jwtKey := []byte(config.LoadEnv("JWT_REFRESH_KEY"))
	expirationTime := time.Now().Add(1 * 1 * time.Minute) // 7 hari

	claims := &model.JWTClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ResetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     TokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Berhasil logout"))
}

func Authenticate(next http.HandlerFunc, accessType int) http.HandlerFunc {
	log.Println("function authenticate dipanggil")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := ValidateUserToken(r, accessType)
		if !isValidToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("gagal authenticate")
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func ValidateUserToken(r *http.Request, accessType int) bool {
	log.Println("function validate user token dipanggil")
	isAccessTokenValid, _, _, userType := ValidateTokenFormCookies(r)
	// isAccessTokenValid, id, username, userType := ValidateTokenFormCookies(r)
	// fmt.Println(id, username, userType, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func ValidateTokenFormCookies(r *http.Request) (bool, int, string, int) {
	log.Println("function validate token form cookies dipanggil")
	jwtKey := []byte(config.LoadEnv("JWT_KEY"))
	cookie, err1 := r.Cookie(TokenName)
	// log.Println(cookie)
	if err1 == nil {
		accessToken := cookie.Value
		accessClaims := &model.JWTClaims{}
		parsedToken, err2 := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err2 == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Username, accessClaims.UserType
		} else {
			log.Println(err2)
		}
	} else {
		log.Println(err1)
	}
	return false, -1, "", -1
}

func GetIdFromCookie(r *http.Request) int {
	jwtKey := []byte(config.LoadEnv("JWT_KEY"))
	cookie, err1 := r.Cookie(TokenName)
	if err1 == nil {
		accessToken := cookie.Value
		accessClaims := &model.JWTClaims{}
		parsedToken, err2 := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err2 == nil && parsedToken.Valid {
			return accessClaims.ID
		} else {
			log.Println(err2)
		}
	} else {
		log.Println(err1)
	}
	return -1
}