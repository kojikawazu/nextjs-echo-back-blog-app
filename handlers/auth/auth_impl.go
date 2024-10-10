package handlers_auth

import (
	services_auth "backend/services/auth"
	services_users "backend/services/users"
)

type AuthHandler struct {
	UserService services_users.UserService
	AuthService services_auth.AuthService
}

// コンストラクタ
func NewAuthHandler(userService services_users.UserService, authService services_auth.AuthService) *AuthHandler {
	return &AuthHandler{
		UserService: userService,
		AuthService: authService,
	}
}
