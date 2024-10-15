package services_blogs

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockBlogService struct {
	mock.Mock
}

func (m *MockBlogService) FetchBlogs() ([]models.BlogData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogService) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BlogData), args.Error(1)
}

func (m *MockBlogService) FetchBlogById(id string) (*models.BlogData, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

func (m *MockBlogService) CreateBlog(userId, title, githubUrl, category, description, tags string) (models.BlogData, error) {
	args := m.Called(userId, title, githubUrl, category, description, tags)
	if args.Get(0) == nil {
		return models.BlogData{}, args.Error(1)
	}
	return args.Get(0).(models.BlogData), args.Error(1)
}

func (m *MockBlogService) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(id, title, githubUrl, category, description, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

func (m *MockBlogService) DeleteBlog(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
