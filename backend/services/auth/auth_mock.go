package services_auth

import "github.com/stretchr/testify/mock"

type MockAuthService struct {
	mock.Mock
}

// Login メソッドのモック実装
func (m *MockAuthService) Login(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}
