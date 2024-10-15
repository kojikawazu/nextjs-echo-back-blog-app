package repositories_users

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	args := m.Called(email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FetchUserById(id string) (*models.UserData, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserData), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(id, name, email, password string) (*models.UserData, error) {
	args := m.Called(id, name, email, password)
	if args.Get(0) != nil {
		return args.Get(0).(*models.UserData), args.Error(1)
	}
	return nil, args.Error(1)
}
