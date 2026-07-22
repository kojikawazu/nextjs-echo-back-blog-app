package repositories_blog_comments

import "backend/models"

// CommentRepository はコメントのデータアクセスを担うリポジトリインターフェース。
type CommentRepository interface {
	// FetchCommentsByBlogId はブログIDに一致するコメント情報を取得する。
	FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error)
	// CreateComment はコメント情報を新規作成する。
	CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error)
}

// CommentRepositoryImpl は CommentRepository の実装。
type CommentRepositoryImpl struct{}

// NewCommentRepository は CommentRepository を実装した CommentRepositoryImpl のポインタを生成する。
func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}
