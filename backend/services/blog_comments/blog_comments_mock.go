package services_blog_comments

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockCommentService は CommentService のテスト用モック。
type MockCommentService struct {
	mock.Mock
}

// FetchCommentsByBlogId は CommentService.FetchCommentsByBlogId のテスト用モック。
//
// 引数:
//   - blogId: 取得対象のブログID
//
// 戻り値:
//   - []models.BlogCommentsData: モックの戻り値設定に従うコメントデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockCommentService) FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error) {
	args := m.Called(blogId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BlogCommentsData), args.Error(1)
}

// CreateComment は CommentService.CreateComment のテスト用モック。
//
// 引数:
//   - blogId: コメント対象のブログID
//   - guestUser: コメントを投稿するゲストユーザー名
//   - comment: コメント本文
//
// 戻り値:
//   - *models.BlogCommentsData: モックの戻り値設定に従う作成済みコメントデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockCommentService) CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error) {
	args := m.Called(blogId, guestUser, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogCommentsData), args.Error(1)
}
