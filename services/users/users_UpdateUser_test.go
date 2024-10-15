package services_users

import (
	"backend/models"
	repositories_users "backend/repositories/users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_UpdateUser(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUser := &models.UserData{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockUserRepository.On("UpdateUser", "1", "John Doe", "john@example.com", "123").Return(mockUser, nil)

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("1", "John Doe", "john@example.com", "123")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}

func TestService_UpdateUser_InvalidId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("", "John Doe", "john@example.com", "123")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "id is required", err.Error())

	// モックが期待通りに呼び出されていないかを確認
	mockUserRepository.AssertNotCalled(t, "UpdateUser", "", "John Doe", "john@example.com", "123")
}

func TestService_UpdateUser_InvalidName(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("123", "", "john@example.com", "123")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "name is required", err.Error())

	// モックが期待通りに呼び出されていないかを確認
	mockUserRepository.AssertNotCalled(t, "UpdateUser", "123", "", "john@example.com", "123")
}

func TestService_UpdateUser_InvalidEmail01(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("123", "John Doe", "", "123")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "email is required", err.Error())

	// モックが期待通りに呼び出されていないかを確認
	mockUserRepository.AssertNotCalled(t, "UpdateUser", "123", "John Doe", "", "123")
}

func TestService_UpdateUser_InvalidEmail02(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("123", "John Doe", "aaa", "123")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "invalid email format", err.Error())

	// モックが期待通りに呼び出されていないかを確認
	mockUserRepository.AssertNotCalled(t, "UpdateUser", "123", "John Doe", "aaa", "123")
}

func TestService_UpdateUser_InvalidPassword(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("123", "John Doe", "john@example.com", "")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "password is required", err.Error())

	// モックが期待通りに呼び出されていないかを確認
	mockUserRepository.AssertNotCalled(t, "UpdateUser", "123", "John Doe", "john@example.com", "")
}

func TestService_UpdateUser_NotUpdate(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUserRepository.On("UpdateUser", "123", "John Doe", "john@example.com", "123").Return(nil, errors.New("failed to update user"))

	// サービス層メソッドの実行
	user, err := userService.UpdateUser("123", "John Doe", "john@example.com", "123")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "failed to update user", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}
