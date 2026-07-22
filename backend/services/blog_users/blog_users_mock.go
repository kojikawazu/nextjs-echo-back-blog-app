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
//
// 引数:
//   - email: 取得対象のメールアドレス
//   - password: 取得対象のパスワード
//
// 戻り値:
//   - *models.BlogUsersData: モックの戻り値設定に従うユーザーデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockUserService) FetchUserByEmailAndPassword(email, password string) (*models.BlogUsersData, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogUsersData), args.Error(1)
}

// FetchUserById は UserService.FetchUserById のテスト用モック。
//
// 引数:
//   - id: 取得対象のユーザーID
//
// 戻り値:
//   - *models.BlogUsersData: モックの戻り値設定に従うユーザーデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockUserService) FetchUserById(id string) (*models.BlogUsersData, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogUsersData), args.Error(1)
}

// UpdateUser は UserService.UpdateUser のテスト用モック。
//
// 引数:
//   - id: 更新対象のユーザーID
//   - name: 更新後のユーザー名
//   - email: 更新後のメールアドレス
//   - password: 本人確認用の現在のパスワード
//   - newPassword: 更新後の新しいパスワード
//
// 戻り値:
//   - *models.BlogUsersData: モックの戻り値設定に従う更新後ユーザーデータ
//   - error: モックの戻り値設定に従うエラー
func (m *MockUserService) UpdateUser(id, name, email, password, newPassword string) (*models.BlogUsersData, error) {
	args := m.Called(id, name, email, password, newPassword)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BlogUsersData), args.Error(1)
}
