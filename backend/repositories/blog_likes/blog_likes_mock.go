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
func (m *MockBlogLikeRepository) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error) {
	args := m.Called(visitId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogLikesData), args.Error(1)
	}
	return nil, args.Error(1)
}

// IsBlogLiked は BlogLikeRepository.IsBlogLiked のモック実装。
func (m *MockBlogLikeRepository) IsBlogLiked(blogId, visitId string) (bool, error) {
	args := m.Called(blogId, visitId)
	return args.Bool(0), args.Error(1)
}

// CreateBlogLike は BlogLikeRepository.CreateBlogLike のモック実装。
func (m *MockBlogLikeRepository) CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error) {
	args := m.Called(blogId, visitId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogLikesData), args.Error(1)
}

// DeleteBlogLike は BlogLikeRepository.DeleteBlogLike のモック実装。
func (m *MockBlogLikeRepository) DeleteBlogLike(blogId, visitId string) error {
	args := m.Called(blogId, visitId)
	return args.Error(0)
}
