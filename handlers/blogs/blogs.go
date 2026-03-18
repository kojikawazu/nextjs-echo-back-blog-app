package handlers_blogs

import (
	services_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
)

type BlogHandler struct {
	BlogService services_blogs.BlogService
	CookieUtils utils_cookie.CookieUtils
}

// コンストラクタ
func NewBlogHandler(blogService services_blogs.BlogService, cookieUtils utils_cookie.CookieUtils) *BlogHandler {
	return &BlogHandler{
		BlogService: blogService,
		CookieUtils: cookieUtils,
	}
}
