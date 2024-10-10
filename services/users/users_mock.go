package services_users

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserData), args.Error(1)
}
