package repositories_blog_comments

import "backend/models"

// CommentRepository はコメントのデータアクセスを担うリポジトリインターフェース。
type CommentRepository interface {
	// FetchCommentsByBlogId はブログIDに一致するコメント情報を取得する。
	//
	// 引数:
	//   - blogId: 取得対象のブログID
	//
	// 戻り値:
	//   - []models.BlogCommentsData: 取得したコメント情報の一覧
	//   - error: 取得に失敗した場合のエラー
	FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error)
	// CreateComment はコメント情報を新規作成する。
	//
	// 引数:
	//   - blogId: コメント対象のブログID
	//   - guestUser: コメントを投稿したゲストユーザー名
	//   - comment: コメント本文
	//
	// 戻り値:
	//   - *models.BlogCommentsData: 作成したコメント情報
	//   - error: 作成に失敗した場合のエラー
	CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error)
}

// CommentRepositoryImpl は CommentRepository の実装。
type CommentRepositoryImpl struct{}

// NewCommentRepository は CommentRepository を実装した CommentRepositoryImpl のポインタを生成する。
//
// 戻り値:
//   - CommentRepository: 生成したリポジトリ実装
func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}
