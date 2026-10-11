package services_blog_users

import (
	"backend/models"
	repositories_blog_users "backend/repositories/blog_users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: 保存されている bcrypt ハッシュと一致すればユーザーを返す。
func TestService_FetchUserByEmailAndPassword(t *testing.T) {
	mockUserRepository := new(repositories_blog_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	mockUser := &models.BlogUsersData{
		ID:       "1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: mustHash(t, "password123"),
	}
	mockUserRepository.On("FetchBlogUsersByEmail", "john@example.com").Return(mockUser, nil)

	user, err := userService.FetchUserByEmailAndPassword("john@example.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	mockUserRepository.AssertExpectations(t)
}

// errRepositoryFailure はリポジトリの想定外エラー（DB 障害等）を表すテスト用エラー。
var errRepositoryFailure = errors.New("connection refused")

// 異常系: DB 障害を認証失敗（ErrInvalidCredentials = 401）に丸めない。
func TestService_FetchUserByEmailAndPassword_RepositoryErrorIsNotInvalidCredentials(t *testing.T) {
	mockUserRepository := new(repositories_blog_users.MockUserRepository)
	mockUserRepository.On("FetchBlogUsersByEmail", "john@example.com").Return(nil, errRepositoryFailure)

	user, err := NewUserService(mockUserRepository).FetchUserByEmailAndPassword("john@example.com", "password123")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, errRepositoryFailure)
	assert.NotErrorIs(t, err, ErrInvalidCredentials)
	mockUserRepository.AssertExpectations(t)
}

// TestService_FetchUserByEmailAndPassword_InvalidCases は入力不正・認証失敗時の振る舞いを検証する。
func TestService_FetchUserByEmailAndPassword_InvalidCases(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		setup    func(m *repositories_blog_users.MockUserRepository)
		wantErr  error
	}{
		{
			// 準正常系: 必須項目の欠落
			name: "empty_email_and_password", email: "", password: "",
			setup:   func(m *repositories_blog_users.MockUserRepository) {},
			wantErr: ErrEmailAndPasswordRequired,
		},
		{
			// 準正常系: メール形式の不正
			name: "invalid_email_format", email: "invalid-email", password: "password123",
			setup:   func(m *repositories_blog_users.MockUserRepository) {},
			wantErr: ErrInvalidEmailFormat,
		},
		{
			// 準正常系: ユーザー不在（#137: リポジトリが返す ErrUserNotFound を判定できる）
			name: "user_not_found", email: "john@example.com", password: "password123",
			setup: func(m *repositories_blog_users.MockUserRepository) {
				m.On("FetchBlogUsersByEmail", "john@example.com").Return(nil, repositories_blog_users.ErrUserNotFound)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			// 準正常系: パスワード不一致（ユーザー不在と同じエラーで区別させない）
			name: "wrong_password", email: "john@example.com", password: "wrong-password",
			setup: func(m *repositories_blog_users.MockUserRepository) {
				m.On("FetchBlogUsersByEmail", "john@example.com").Return(&models.BlogUsersData{ID: "1", Password: mustHash(t, "password123")}, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			// 異常系: 保存値が平文のまま（移行漏れ）でも認証を通さない
			name: "plaintext_stored_password", email: "john@example.com", password: "password123",
			setup: func(m *repositories_blog_users.MockUserRepository) {
				m.On("FetchBlogUsersByEmail", "john@example.com").Return(&models.BlogUsersData{ID: "1", Password: "password123"}, nil)
			},
			wantErr: ErrInvalidCredentials,
		},
		{
			// 異常系: リポジトリの想定外エラーは原因を辿れる形で返す（認証失敗と区別する）
			name: "repository_error", email: "john@example.com", password: "password123",
			setup: func(m *repositories_blog_users.MockUserRepository) {
				m.On("FetchBlogUsersByEmail", "john@example.com").Return(nil, errRepositoryFailure)
			},
			wantErr: errRepositoryFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(repositories_blog_users.MockUserRepository)
			tt.setup(mockUserRepository)

			user, err := NewUserService(mockUserRepository).FetchUserByEmailAndPassword(tt.email, tt.password)

			assert.Nil(t, user)
			assert.ErrorIs(t, err, tt.wantErr)
			mockUserRepository.AssertExpectations(t)
		})
	}
}
