package repositories_blogs_likes

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockBlogLikeRepository struct {
	mock.Mock
}

func (m *MockBlogLikeRepository) IsBlogLiked(blogId, visitId string) (bool, error) {
	args := m.Called(blogId, visitId)
	return args.Bool(0), args.Error(1)
}

func (m *MockBlogLikeRepository) CreateBlogLike(blogId, visitId string) (*models.BlogLikeData, error) {
	args := m.Called(blogId, visitId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogLikeData), args.Error(1)
}

func (m *MockBlogLikeRepository) DeleteBlogLike(blogId, visitId string) error {
	args := m.Called(blogId, visitId)
	return args.Error(0)
}
