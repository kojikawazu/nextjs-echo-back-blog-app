package models

import "github.com/golang-jwt/jwt"

// ユーザー情報のペイロード
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}
