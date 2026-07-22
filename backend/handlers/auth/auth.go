package handlers_auth

import (
	services_auth "backend/services/auth"
	services_users "backend/services/blog_users"
)

// AuthHandler は認証系エンドポイントのハンドラ。
type AuthHandler struct {
	UserService services_users.UserService
	AuthService services_auth.AuthService
}

// NewAuthHandler は AuthHandler を生成する。
//
// 引数:
//   - userService: ユーザーデータ取得に用いるユーザーサービス
//   - authService: 認証処理に用いる認証サービス
//
// 戻り値:
//   - *AuthHandler: 生成した認証ハンドラ
func NewAuthHandler(userService services_users.UserService, authService services_auth.AuthService) *AuthHandler {
	return &AuthHandler{
		UserService: userService,
		AuthService: authService,
	}
}
