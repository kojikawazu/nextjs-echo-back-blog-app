package models

import "github.com/golang-jwt/jwt"

// ユーザー情報のペイロード
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 訪問者IDのペイロード
type ClaimsVisitId struct {
	VisitId string `json:"visit_id"`
	jwt.StandardClaims
}
