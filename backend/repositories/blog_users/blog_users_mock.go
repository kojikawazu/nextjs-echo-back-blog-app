package repositories_blog_users

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FetchBlogUsersByEmailAndPassword(email, password string) (*models.BlogUsersData, error) {
	args := m.Called(email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogUsersData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FetchBlogUsersById(id string) (*models.BlogUsersData, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogUsersData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateBlogUsers(id, name, email, password string) (*models.BlogUsersData, error) {
	args := m.Called(id, name, email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogUsersData), args.Error(1)
	}
	return nil, args.Error(1)
}
