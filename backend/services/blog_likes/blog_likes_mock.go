package services_blog_likes

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockBlogLikeService は BlogLikeService のテスト用モック。
type MockBlogLikeService struct {
	mock.Mock
}

// FetchBlogLikesByVisitId は BlogLikeService.FetchBlogLikesByVisitId のテスト用モック。
//
// 引数:
//   - visitId: 取得対象の訪問者ID
//
// 戻り値:
//   - []models.BlogLikesData: モックの戻り値設定に従ういいねデータの一覧
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeService) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error) {
	args := m.Called(visitId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogLikesData), args.Error(1)
	}
	return nil, args.Error(1)
}

// IsBlogLiked は BlogLikeService.IsBlogLiked のテスト用モック。
//
// 引数:
//   - blogId: 確認対象のブログID
//   - visitId: 確認対象の訪問者ID
//
// 戻り値:
//   - bool: モックの戻り値設定に従ういいね存在有無
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeService) IsBlogLiked(blogId, visitId string) (bool, error) {
	args := m.Called(blogId, visitId)
	return args.Bool(0), args.Error(1)
}

// CreateBlogLike は BlogLikeService.CreateBlogLike のテスト用モック。
//
// 引数:
//   - blogId: いいね対象のブログID
//   - visitId: いいねする訪問者ID
//
// 戻り値:
//   - *models.BlogLikesData: モックの戻り値設定に従う作成済みいいねデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeService) CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error) {
	args := m.Called(blogId, visitId)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogLikesData), args.Error(1)
	}
	return nil, args.Error(1)
}

// DeleteBlogLike は BlogLikeService.DeleteBlogLike のテスト用モック。
//
// 引数:
//   - blogId: 削除対象のブログID
//   - visitId: 削除対象の訪問者ID
//
// 戻り値:
//   - error: モックの戻り値設定に従うエラー
func (m *MockBlogLikeService) DeleteBlogLike(blogId, visitId string) error {
	args := m.Called(blogId, visitId)
	return args.Error(0)
}
