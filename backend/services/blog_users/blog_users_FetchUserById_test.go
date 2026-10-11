package services_blog_users

import (
	"backend/models"
	repositories_users "backend/repositories/blog_users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchUserById(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUser := &models.BlogUsersData{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockUserRepository.On("FetchBlogUsersById", "1").Return(mockUser, nil)

	// サービス層メソッドの実行
	user, err := userService.FetchUserById("1")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

// TestService_FetchUserById_Errors は FetchUserById の入力不正・取得失敗時の振る舞いを検証する。
func TestService_FetchUserById_Errors(t *testing.T) {
	errConnection := errors.New("connection refused")

	tests := []struct {
		name         string
		id           string
		setup        func(m *repositories_users.MockUserRepository)
		wantErr      error
		wantNotErr   error
		expectCalled bool
	}{
		{
			// 準正常系: id が空の場合はリポジトリを呼ばずに ErrIDRequired を返す
			name:    "empty_id",
			id:      "",
			setup:   func(m *repositories_users.MockUserRepository) {},
			wantErr: ErrIDRequired,
		},
		{
			// 準正常系: ユーザー不在はリポジトリの ErrUserNotFound をサービスの ErrUserNotFound（404）に変換する
			name: "user_not_found",
			id:   "123",
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "123").Return(nil, repositories_users.ErrUserNotFound)
			},
			wantErr:      ErrUserNotFound,
			expectCalled: true,
		},
		{
			// 異常系: リポジトリの想定外エラーは原因を辿れる形で返し、不在（404）には丸めない
			name: "repository_error",
			id:   "123",
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "123").Return(nil, errConnection)
			},
			wantErr:      errConnection,
			wantNotErr:   ErrUserNotFound,
			expectCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(repositories_users.MockUserRepository)
			tt.setup(mockUserRepository)

			user, err := NewUserService(mockUserRepository).FetchUserById(tt.id)

			assert.Nil(t, user)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantNotErr != nil {
				assert.NotErrorIs(t, err, tt.wantNotErr)
			}
			if tt.expectCalled {
				mockUserRepository.AssertExpectations(t)
			} else {
				mockUserRepository.AssertNotCalled(t, "FetchBlogUsersById", tt.id)
			}
		})
	}
}
