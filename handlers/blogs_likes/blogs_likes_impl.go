package handlers_blogs_likes

import (
	services_blogs_likes "backend/services/blogs_likes"
	utils_cookie "backend/utils/cookie"
)

type BlogLikeHandler struct {
	BlogLikeService services_blogs_likes.BlogLikeService
	CookieUtils     utils_cookie.CookieUtils
}

// コンストラクタ
func NewBlogLikeHandler(blogLikeService services_blogs_likes.BlogLikeService, cookieUtils utils_cookie.CookieUtils) *BlogLikeHandler {
	return &BlogLikeHandler{
		BlogLikeService: blogLikeService,
		CookieUtils:     cookieUtils,
	}
}
