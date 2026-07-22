package services_auth

import "github.com/stretchr/testify/mock"

// MockAuthService は AuthService のテスト用モック。
type MockAuthService struct {
	mock.Mock
}

// Login は AuthService.Login のテスト用モック。
//
// 引数:
//   - email: 検証対象のメールアドレス
//   - password: 検証対象のパスワード
//
// 戻り値:
//   - error: モックの戻り値設定に従うエラー
func (m *MockAuthService) Login(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}
