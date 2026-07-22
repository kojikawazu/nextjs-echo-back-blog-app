package services_blogs

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockBlogService は BlogService のテスト用モック。
type MockBlogService struct {
	mock.Mock
}

// FetchBlogs は BlogService.FetchBlogs のテスト用モック。
func (m *MockBlogService) FetchBlogs() ([]models.BlogData, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogsByUserId は BlogService.FetchBlogsByUserId のテスト用モック。
func (m *MockBlogService) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BlogData), args.Error(1)
}

// FetchBlogById は BlogService.FetchBlogById のテスト用モック。
func (m *MockBlogService) FetchBlogById(id string) (*models.BlogData, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

// CreateBlog は BlogService.CreateBlog のテスト用モック。
func (m *MockBlogService) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(userId, title, githubUrl, category, description, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

// UpdateBlog は BlogService.UpdateBlog のテスト用モック。
func (m *MockBlogService) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	args := m.Called(id, title, githubUrl, category, description, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogData), args.Error(1)
}

// DeleteBlog は BlogService.DeleteBlog のテスト用モック。
func (m *MockBlogService) DeleteBlog(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// FetchBlogCategories は BlogService.FetchBlogCategories のテスト用モック。
func (m *MockBlogService) FetchBlogCategories() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogTags は BlogService.FetchBlogTags のテスト用モック。
func (m *MockBlogService) FetchBlogTags() ([]string, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogPopular は BlogService.FetchBlogPopular のテスト用モック。
func (m *MockBlogService) FetchBlogPopular(count int) ([]models.BlogData, error) {
	args := m.Called(count)
	if args.Get(0) != nil {
		return args.Get(0).([]models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}
