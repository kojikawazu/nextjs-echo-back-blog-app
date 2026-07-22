package repositories_blog_likes

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockBlogLikeRepository は BlogLikeRepository のテスト用モック。
type MockBlogLikeRepository struct {
	mock.Mock
}

// FetchBlogLikesByVisitId は BlogLikeRepository.FetchBlogLikesByVisitId のモック実装。
//
// 引数:
//   - visitId: 取得対象の訪問者ID
//
// 戻り値:
//   - []models.BlogLikesData: モックの戻り値設定に従ういいねデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeRepository) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error) {
	args := m.Called(visitId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogLikesData), args.Error(1)
	}
	return nil, args.Error(1)
}

// IsBlogLiked は BlogLikeRepository.IsBlogLiked のモック実装。
//
// 引数:
//   - blogId: 確認対象のブログID
//   - visitId: 確認対象の訪問者ID
//
// 戻り値:
//   - bool: モックの戻り値設定に従ういいねの存在有無
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeRepository) IsBlogLiked(blogId, visitId string) (bool, error) {
	args := m.Called(blogId, visitId)
	return args.Bool(0), args.Error(1)
}

// CreateBlogLike は BlogLikeRepository.CreateBlogLike のモック実装。
//
// 引数:
//   - blogId: いいね対象のブログID
//   - visitId: いいねした訪問者ID
//
// 戻り値:
//   - *models.BlogLikesData: モックの戻り値設定に従う作成後のいいねデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeRepository) CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error) {
	args := m.Called(blogId, visitId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogLikesData), args.Error(1)
}

// DeleteBlogLike は BlogLikeRepository.DeleteBlogLike のモック実装。
//
// 引数:
//   - blogId: 削除対象のブログID
//   - visitId: 削除対象の訪問者ID
//
// 戻り値:
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeRepository) DeleteBlogLike(blogId, visitId string) error {
	args := m.Called(blogId, visitId)
	return args.Error(0)
}
