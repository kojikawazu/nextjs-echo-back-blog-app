package services_blog_likes

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockBlogLikeService struct {
	mock.Mock
}

func (m *MockBlogLikeService) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error) {
	args := m.Called(visitId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogLikesData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogLikeService) IsBlogLiked(blogId, visitId string) (bool, error) {
	args := m.Called(blogId, visitId)
	return args.Bool(0), args.Error(1)
}

func (m *MockBlogLikeService) CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error) {
	args := m.Called(blogId, visitId)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogLikesData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogLikeService) DeleteBlogLike(blogId, visitId string) error {
	args := m.Called(blogId, visitId)
	return args.Error(0)
}
