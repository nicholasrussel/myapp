package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	UserType int    `json:"user_type"`
	jwt.RegisteredClaims
}