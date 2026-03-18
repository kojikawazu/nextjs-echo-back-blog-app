package services_blog_users

import (
	"backend/models"
	repositories_blog_users "backend/repositories/blog_users"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchUserByEmailAndPassword(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_blog_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUser := &models.BlogUsersData{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockUserRepository.On("FetchBlogUsersByEmailAndPassword", "john@example.com", "password123").Return(mockUser, nil)

	// サービス層メソッドの実行
	user, err := userService.FetchUserByEmailAndPassword("john@example.com", "password123")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_FetchUserByEmailAndPassword_InvalidCases(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_blog_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// 1. メールアドレスとパスワードが空の場合
	_, err := userService.FetchUserByEmailAndPassword("", "")
	assert.Error(t, err)
	assert.Equal(t, "email and password are required", err.Error())

	// 2. メールアドレスの形式が無効な場合
	_, err = userService.FetchUserByEmailAndPassword("invalid-email", "password123")
	assert.Error(t, err)
	assert.Equal(t, "invalid email format", err.Error())

	// 3. ユーザーが見つからない場合
	mockUserRepository.On("FetchBlogUsersByEmailAndPassword", "john@example.com", "password123").Return(nil, sql.ErrNoRows)

	_, err = userService.FetchUserByEmailAndPassword("john@example.com", "password123")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}
