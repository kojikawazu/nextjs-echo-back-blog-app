package repositories_blogs

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockBlogRepository struct {
	mock.Mock
}

func (m *MockBlogRepository) FetchBlogs() ([]models.BlogData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepository) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	args := m.Called(userId)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepository) FetchBlogById(id string) (*models.BlogData, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepository) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(userId, title, githubUrl, category, description, tags)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepository) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(id, title, githubUrl, category, description, tags)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepository) DeleteBlog(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBlogRepository) FetchBlogCategories() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepository) FetchBlogTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBlogRepository) FetchBlogPopular(count int) ([]models.BlogData, error) {
	args := m.Called(count)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}
