package services_auth

import "github.com/stretchr/testify/mock"

// MockAuthService は AuthService のテスト用モック。
type MockAuthService struct {
	mock.Mock
}

// Login は AuthService.Login のテスト用モック。
func (m *MockAuthService) Login(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}
