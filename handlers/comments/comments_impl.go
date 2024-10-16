package handlers_comments

import services_comments "backend/services/comments"

type CommentHandler struct {
	CommentService services_comments.CommentService
}

// コンストラクタ
func NewCommentHandler(commentService services_comments.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
	}
}
