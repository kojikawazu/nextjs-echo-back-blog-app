package handlers_blog_comments

import services_comments "backend/services/blog_comments"

// CommentHandler はブログコメント系エンドポイントのハンドラ。
type CommentHandler struct {
	CommentService services_comments.CommentService
}

// NewCommentHandler は CommentHandler を生成する。
func NewCommentHandler(commentService services_comments.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
	}
}
