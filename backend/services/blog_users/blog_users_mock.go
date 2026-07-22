package services_blog_users

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockUserService は UserService のテスト用モック。
type MockUserService struct {
	mock.Mock
}

// FetchUserByEmailAndPassword は UserService.FetchUserByEmailAndPassword のテスト用モック。
func (m *MockUserService) FetchUserByEmailAndPassword(email, password string) (*models.BlogUsersData, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogUsersData), args.Error(1)
}

// FetchUserById は UserService.FetchUserById のテスト用モック。
func (m *MockUserService) FetchUserById(id string) (*models.BlogUsersData, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogUsersData), args.Error(1)
}

// UpdateUser は UserService.UpdateUser のテスト用モック。
func (m *MockUserService) UpdateUser(id, name, email, password, newPassword string) (*models.BlogUsersData, error) {
	args := m.Called(id, name, email, password, newPassword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogUsersData), args.Error(1)
}
