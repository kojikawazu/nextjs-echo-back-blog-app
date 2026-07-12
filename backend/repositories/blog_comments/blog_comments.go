package repositories_blog_comments

import "backend/models"

// CommentRepositoryインターフェース
type CommentRepository interface {
	FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error)
	CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error)
}

type CommentRepositoryImpl struct{}

// CommentRepositoryインターフェースを実装したCommentRepositoryImplのポインタを返す
func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}
