package services_blog_comments

import (
	"backend/models"
	repositories_comments "backend/repositories/blog_comments"
)

// CommentService はコメントに関する処理を提供するサービスインターフェース。
type CommentService interface {
	// FetchCommentsByBlogId は指定されたブログIDに一致するコメントデータを取得する。
	FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error)
	// CreateComment はコメントデータを新規作成する。
	CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error)
}

// CommentServiceImpl は CommentService の実装。
type CommentServiceImpl struct {
	CommentRepository repositories_comments.CommentRepository
}

// NewCommentService は CommentService インターフェースを実装した CommentServiceImpl を生成する。
func NewCommentService(
	commentRepository repositories_comments.CommentRepository,
) CommentService {
	return &CommentServiceImpl{
		CommentRepository: commentRepository,
	}
}
