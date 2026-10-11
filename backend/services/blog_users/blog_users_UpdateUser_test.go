package services_blog_users

import (
	"backend/models"
	repositories_users "backend/repositories/blog_users"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_UpdateUser(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUser := &models.BlogUsersData{
		ID:       "1",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: mustHash(t, "123"),
	}
	mockUserRepository.On("FetchBlogUsersById", "1").Return(mockUser, nil)
	mockUserRepository.On("UpdateBlogUsers", "1", "John Doe", "john@example.com", mock.MatchedBy(isHashOf("1234"))).Return(mockUser, nil)

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("1", "John Doe", "john@example.com", "123", "1234")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

// TestService_UpdateUser_InvalidInput は入力不正時にリポジトリを呼ばずに sentinel を返すことを検証する。
func TestService_UpdateUser_InvalidInput(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		userName    string
		email       string
		password    string
		newPassword string
		wantErr     error
	}{
		{
			// 準正常系: id が空
			name: "empty_id", id: "", userName: "John Doe", email: "john@example.com", password: "123", newPassword: "1234",
			wantErr: ErrIDRequired,
		},
		{
			// 準正常系: name が空
			name: "empty_name", id: "123", userName: "", email: "john@example.com", password: "123", newPassword: "1234",
			wantErr: ErrNameRequired,
		},
		{
			// 準正常系: email が空
			name: "empty_email", id: "123", userName: "John Doe", email: "", password: "123", newPassword: "1234",
			wantErr: ErrEmailRequired,
		},
		{
			// 準正常系: email の形式が不正
			name: "invalid_email_format", id: "123", userName: "John Doe", email: "aaa", password: "123", newPassword: "1234",
			wantErr: ErrInvalidEmailFormat,
		},
		{
			// 準正常系: 現在のパスワードが空
			name: "empty_password", id: "123", userName: "John Doe", email: "john@example.com", password: "", newPassword: "1234",
			wantErr: ErrPasswordRequired,
		},
		{
			// 準正常系: 新しいパスワードが空
			name: "empty_new_password", id: "123", userName: "John Doe", email: "john@example.com", password: "123", newPassword: "",
			wantErr: ErrNewPasswordRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(repositories_users.MockUserRepository)

			user, err := NewUserService(mockUserRepository).UpdateUser(tt.id, tt.userName, tt.email, tt.password, tt.newPassword)

			assert.Nil(t, user)
			assert.ErrorIs(t, err, tt.wantErr)
			mockUserRepository.AssertNotCalled(t, "FetchBlogUsersById", mock.Anything)
			mockUserRepository.AssertNotCalled(t, "UpdateBlogUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

// TestService_UpdateUser_RepositoryErrors は本人確認・更新でリポジトリが失敗した場合の振る舞いを検証する。
func TestService_UpdateUser_RepositoryErrors(t *testing.T) {
	errConnection := errors.New("connection refused")
	currentUser := func() *models.BlogUsersData {
		return &models.BlogUsersData{ID: "1", Name: "John Doe", Email: "john@example.com", Password: mustHash(t, "123")}
	}

	tests := []struct {
		name        string
		password    string
		newPassword string
		setup       func(m *repositories_users.MockUserRepository)
		wantErr     error
		wantNotErr  error
		wantUpdate  bool
	}{
		{
			// 準正常系: 本人確認時にユーザー不在 → ErrUserNotFound（404）
			name: "fetch_user_not_found", password: "123", newPassword: "1234",
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "1").Return(nil, repositories_users.ErrUserNotFound)
			},
			wantErr: ErrUserNotFound,
		},
		{
			// 異常系: 本人確認時の想定外エラーは原因を辿れる形で返し、不在（404）には丸めない
			name: "fetch_repository_error", password: "123", newPassword: "1234",
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "1").Return(nil, errConnection)
			},
			wantErr:    errConnection,
			wantNotErr: ErrUserNotFound,
		},
		{
			// 準正常系: 現在のパスワード不一致 → ErrInvalidCurrentPassword（更新しない）
			name: "invalid_current_password", password: "12355", newPassword: "1234",
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "1").Return(currentUser(), nil)
			},
			wantErr: ErrInvalidCurrentPassword,
		},
		{
			// 準正常系: 新しいパスワードが 72 バイトを超える場合は保存せず ErrPasswordTooLong を返す
			name: "new_password_too_long", password: "123", newPassword: strings.Repeat("a", 73),
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "1").Return(currentUser(), nil)
			},
			wantErr: ErrPasswordTooLong,
		},
		{
			// 準正常系: 更新時にユーザー不在（本人確認後に削除された） → ErrUserNotFound（404）
			name: "update_user_not_found", password: "123", newPassword: "1234",
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "1").Return(currentUser(), nil)
				m.On("UpdateBlogUsers", "1", "John Doe", "john@example.com", mock.MatchedBy(isHashOf("1234"))).Return(nil, repositories_users.ErrUserNotFound)
			},
			wantErr:    ErrUserNotFound,
			wantUpdate: true,
		},
		{
			// 異常系: 更新時の想定外エラーは原因を辿れる形で返す
			name: "update_repository_error", password: "123", newPassword: "1234",
			setup: func(m *repositories_users.MockUserRepository) {
				m.On("FetchBlogUsersById", "1").Return(currentUser(), nil)
				m.On("UpdateBlogUsers", "1", "John Doe", "john@example.com", mock.MatchedBy(isHashOf("1234"))).Return(nil, errConnection)
			},
			wantErr:    errConnection,
			wantNotErr: ErrUserNotFound,
			wantUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(repositories_users.MockUserRepository)
			tt.setup(mockUserRepository)

			user, err := NewUserService(mockUserRepository).UpdateUser("1", "John Doe", "john@example.com", tt.password, tt.newPassword)

			assert.Nil(t, user)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantNotErr != nil {
				assert.NotErrorIs(t, err, tt.wantNotErr)
			}
			mockUserRepository.AssertExpectations(t)
			if !tt.wantUpdate {
				mockUserRepository.AssertNotCalled(t, "UpdateBlogUsers", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}
