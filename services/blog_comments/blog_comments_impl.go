package services_blog_comments

import (
	"backend/models"
	repositories_comments "backend/repositories/blog_comments"
)

// CommentServiceインターフェース
type CommentService interface {
	FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error)
	CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error)
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
