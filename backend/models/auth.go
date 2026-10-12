package models

import "github.com/golang-jwt/jwt/v5"

// Claims はユーザー情報を保持する JWT ペイロード。
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ClaimsVisitId は訪問者IDを保持する JWT ペイロード。
type ClaimsVisitId struct {
	VisitId string `json:"visit_id"`
	jwt.RegisteredClaims
}
