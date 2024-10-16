package repositories_comments

import "backend/models"

// CommentRepositoryインターフェース
type CommentRepository interface {
	FetchCommentsByBlogId(blogId string) ([]models.CommentData, error)
	CreateComment(blogId, guestUser, comment string) (*models.CommentData, error)
}

type CommentRepositoryImpl struct{}

// CommentRepositoryインターフェースを実装したCommentRepositoryImplのポインタを返す
func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}
