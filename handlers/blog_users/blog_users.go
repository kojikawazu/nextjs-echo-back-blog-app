package handlers_blog_users

import (
	services_users "backend/services/blog_users"
	utils_cookie "backend/utils/cookie"
)

type BlogUsersHandler struct {
	UserService services_users.UserService
	CookieUtils utils_cookie.CookieUtils
}

// コンストラクタ
func NewBlogUsersHandler(userService services_users.UserService, cookieUtils utils_cookie.CookieUtils) *BlogUsersHandler {
	return &BlogUsersHandler{
		UserService: userService,
		CookieUtils: cookieUtils,
	}
}
