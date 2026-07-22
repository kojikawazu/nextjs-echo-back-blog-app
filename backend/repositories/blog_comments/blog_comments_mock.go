package repositories_blog_comments

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockCommentRepository は CommentRepository のテスト用モック。
type MockCommentRepository struct {
	mock.Mock
}

// FetchCommentsByBlogId は CommentRepository.FetchCommentsByBlogId のモック実装。
//
// 引数:
//   - blogId: 取得対象のブログID
//
// 戻り値:
//   - []models.BlogCommentsData: モックの戻り値設定に従うコメント情報の一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockCommentRepository) FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error) {
	args := m.Called(blogId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogCommentsData), args.Error(1)
	}
	return nil, args.Error(1)
}

// CreateComment は CommentRepository.CreateComment のモック実装。
//
// 引数:
//   - blogId: コメント対象のブログID
//   - guestUser: コメントを投稿したゲストユーザー名
//   - comment: コメント本文
//
// 戻り値:
//   - *models.BlogCommentsData: モックの戻り値設定に従う作成後のコメント情報
//   - error: モックの戻り値設定に従うエラー
func (m *MockCommentRepository) CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error) {
	args := m.Called(blogId, guestUser, comment)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogCommentsData), args.Error(1)
	}
	return nil, args.Error(1)
}
