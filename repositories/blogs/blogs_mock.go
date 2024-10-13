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

func (m *MockBlogRepository) FetchBlogByUserId(userId string) (*models.BlogData, error) {
	args := m.Called(userId)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogData), args.Error(1)
	}
	return nil, args.Error(1)
}
