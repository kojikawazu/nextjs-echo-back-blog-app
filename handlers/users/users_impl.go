package handlers_users

import (
	services_users "backend/services/users"
	utils_cookie "backend/utils/cookie"
)

type UserHandler struct {
	UserService services_users.UserService
	CookieUtils utils_cookie.CookieUtils
}

// コンストラクタ
func NewUserHandler(userService services_users.UserService, cookieUtils utils_cookie.CookieUtils) *UserHandler {
	return &UserHandler{
		UserService: userService,
		CookieUtils: cookieUtils,
	}
}
