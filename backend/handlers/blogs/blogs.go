package handlers_blogs

import (
	services_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
)

// BlogHandler はブログ系エンドポイントのハンドラ。
type BlogHandler struct {
	BlogService services_blogs.BlogService
	CookieUtils utils_cookie.CookieUtils
}

// NewBlogHandler は BlogHandler を生成する。
func NewBlogHandler(blogService services_blogs.BlogService, cookieUtils utils_cookie.CookieUtils) *BlogHandler {
	return &BlogHandler{
		BlogService: blogService,
		CookieUtils: cookieUtils,
	}
}
