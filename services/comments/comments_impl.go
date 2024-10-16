package services_comments

import (
	"backend/models"
	repositories_comments "backend/repositories/comments"
)

// CommentServiceインターフェース
type CommentService interface {
	FetchCommentsByBlogId(blogId string) ([]models.CommentData, error)
	CreateComment(blogId, guestUser, comment string) (*models.CommentData, error)
}

type CommentServiceImpl struct {
	CommentRepository repositories_comments.CommentRepository
}

// CommentServiceインターフェースを実装したCommentServiceImplのポインタを返す
func NewCommentService(
	commentRepository repositories_comments.CommentRepository,
) CommentService {
	return &CommentServiceImpl{
		CommentRepository: commentRepository,
	}
}
