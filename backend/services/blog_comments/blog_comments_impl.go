package services_blog_comments

import (
	"backend/models"
	repositories_comments "backend/repositories/blog_comments"
)

// CommentService はコメントに関する処理を提供するサービスインターフェース。
type CommentService interface {
	// FetchCommentsByBlogId は指定されたブログIDに一致するコメントデータを取得する。
	//
	// 引数:
	//   - blogId: 取得対象のブログID
	//
	// 戻り値:
	//   - []models.BlogCommentsData: 取得したコメントデータの一覧
	//   - error: 取得に失敗した場合のエラー
	FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error)
	// CreateComment はコメントデータを新規作成する。
	//
	// 引数:
	//   - blogId: コメント対象のブログID
	//   - guestUser: コメントを投稿するゲストユーザー名
	//   - comment: コメント本文
	//
	// 戻り値:
	//   - *models.BlogCommentsData: 作成されたコメントデータ
	//   - error: 作成に失敗した場合のエラー
	CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error)
}

// CommentServiceImpl は CommentService の実装。
type CommentServiceImpl struct {
	CommentRepository repositories_comments.CommentRepository
}

// NewCommentService は CommentService インターフェースを実装した CommentServiceImpl を生成する。
//
// 引数:
//   - commentRepository: コメントデータへのアクセスを担うリポジトリ
//
// 戻り値:
//   - CommentService: 生成されたコメントサービス
func NewCommentService(
	commentRepository repositories_comments.CommentRepository,
) CommentService {
	return &CommentServiceImpl{
		CommentRepository: commentRepository,
	}
}
