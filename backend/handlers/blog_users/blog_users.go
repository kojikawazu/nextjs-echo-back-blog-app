package handlers_blog_users

import (
	services_users "backend/services/blog_users"
	utils_cookie "backend/utils/cookie"
)

// BlogUsersHandler はブログユーザー系エンドポイントのハンドラ。
type BlogUsersHandler struct {
	UserService services_users.UserService
	CookieUtils utils_cookie.CookieUtils
}

// NewBlogUsersHandler は BlogUsersHandler を生成する。
//
// 引数:
//   - userService: ユーザーデータの取得/更新に用いるユーザーサービス
//   - cookieUtils: 認証トークンの取得/再発行に用いるクッキーユーティリティ
//
// 戻り値:
//   - *BlogUsersHandler: 生成したユーザーハンドラ
func NewBlogUsersHandler(userService services_users.UserService, cookieUtils utils_cookie.CookieUtils) *BlogUsersHandler {
	return &BlogUsersHandler{
		UserService: userService,
		CookieUtils: cookieUtils,
	}
}
