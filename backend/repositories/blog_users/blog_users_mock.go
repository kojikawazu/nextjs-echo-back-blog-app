package repositories_blog_users

import (
	"backend/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository は BlogUsersRepository のテスト用モック。
type MockUserRepository struct {
	mock.Mock
}

// FetchBlogUsersByEmail は BlogUsersRepository.FetchBlogUsersByEmail のモック実装。
//
// 引数:
//   - email: 取得対象のメールアドレス
//
// 戻り値:
//   - *models.BlogUsersData: モックの戻り値設定に従うユーザー情報
//   - error: モックの戻り値設定に従うエラー
func (m *MockUserRepository) FetchBlogUsersByEmail(email string) (*models.BlogUsersData, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogUsersData), args.Error(1)
	}
	return nil, args.Error(1)
}

// FetchBlogUsersById は BlogUsersRepository.FetchBlogUsersById のモック実装。
//
// 引数:
//   - id: 取得対象のユーザーID
//
// 戻り値:
//   - *models.BlogUsersData: モックの戻り値設定に従うユーザー情報
//   - error: モックの戻り値設定に従うエラー
func (m *MockUserRepository) FetchBlogUsersById(id string) (*models.BlogUsersData, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogUsersData), args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateBlogUsers は BlogUsersRepository.UpdateBlogUsers のモック実装。
//
// 引数:
//   - id: 更新対象のユーザーID
//   - name: 更新後のユーザー名
//   - email: 更新後のメールアドレス
//   - passwordHash: 更新後のパスワードの bcrypt ハッシュ
//
// 戻り値:
//   - *models.BlogUsersData: モックの戻り値設定に従う更新後のユーザー情報
//   - error: モックの戻り値設定に従うエラー
func (m *MockUserRepository) UpdateBlogUsers(id, name, email, passwordHash string) (*models.BlogUsersData, error) {
	args := m.Called(id, name, email, passwordHash)
	if args.Get(0) != nil {
		return args.Get(0).(*models.BlogUsersData), args.Error(1)
	}
	return nil, args.Error(1)
}
