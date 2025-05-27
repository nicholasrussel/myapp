package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nicholasrussel/myapp/config"
	"github.com/nicholasrussel/myapp/internal/model"
)

var TokenName = "loginToken"
var RefreshTokenName = "refreshToken"

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

func ResetUserToken(c *gin.Context) {
	c.SetCookie(TokenName, "", 0, "/", "", false, true)
	c.SetCookie(RefreshTokenName, "", 0, "/", "", false, true)
	c.String(http.StatusOK, "Berhasil logout")
}

func Authenticate(next gin.HandlerFunc, accessType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !ValidateUserToken(c, accessType) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		next(c)
	}
}

func ValidateUserToken(c *gin.Context, accessType int) bool {
	isAccessTokenValid, _, _, userType := ValidateTokenFromCookies(c)
	if isAccessTokenValid && userType == accessType {
		return true
	}
	return false
}

func ValidateTokenFromCookies(c *gin.Context) (bool, int, string, int) {
	jwtKey := []byte(config.LoadEnv("JWT_KEY"))
	accessToken, err := c.Cookie(TokenName)
	if err != nil {
		log.Println(err)
		return false, -1, "", -1
	}

	accessClaims := &model.JWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err == nil && parsedToken.Valid {
		return true, accessClaims.ID, accessClaims.Username, accessClaims.UserType
	}
	log.Println(err)
	return false, -1, "", -1
}

func GetIdFromCookie(c *gin.Context) int {
	jwtKey := []byte(config.LoadEnv("JWT_KEY"))
	accessToken, err := c.Cookie(TokenName)
	if err != nil {
		log.Println(err)
		return -1
	}

	accessClaims := &model.JWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err == nil && parsedToken.Valid {
		return accessClaims.ID
	}
	log.Println(err)
	return -1
}
