package handlers_users

import (
	services_users "backend/services/users"
)

type UserHandler struct {
	UserService services_users.UserService
}

// コンストラクタ
func NewUserHandler(userService services_users.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}
