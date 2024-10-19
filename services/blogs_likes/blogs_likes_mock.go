package services_blogs_likes

import "github.com/stretchr/testify/mock"

type MockBlogLikeService struct {
	mock.Mock
}

func (m *MockBlogLikeService) IsBlogLiked(blogId, visitId string) (bool, error) {
	args := m.Called(blogId, visitId)
	return args.Bool(0), args.Error(1)
}

func (m *MockBlogLikeService) CreateBlogLike(blogId, visitId string) error {
	args := m.Called(blogId, visitId)
	return args.Error(0)
}

func (m *MockBlogLikeService) DeleteBlogLike(blogId, visitId string) error {
	args := m.Called(blogId, visitId)
	return args.Error(0)
}
