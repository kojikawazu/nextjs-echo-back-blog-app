package handlers_blogs

import (
	services_blogs "backend/services/blogs"
)

type BlogHandler struct {
	BlogService services_blogs.BlogService
}

// コンストラクタ
func NewBlogHandler(blogService services_blogs.BlogService) *BlogHandler {
	return &BlogHandler{
		BlogService: blogService,
	}
}
