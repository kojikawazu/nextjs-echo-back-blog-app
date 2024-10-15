package services_users

import (
	"backend/models"
	repositories_users "backend/repositories/users"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchUserById(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUser := &models.UserData{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockUserRepository.On("FetchUserById", "1").Return(mockUser, nil)

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

func TestService_FetchUserById_InvalidId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUserRepository.On("FetchUserById", "").Return(nil, errors.New("user not found"))

	// サービス層メソッドの実行
	user, err := userService.FetchUserById("")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "id is required", err.Error())

	// モックが期待通りに呼びだされていないかを確認
	mockUserRepository.AssertNotCalled(t, "FetchUserById", "")
}

func TestService_FetchUserById_NotUser(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockUserRepository := new(repositories_users.MockUserRepository)
	userService := NewUserService(mockUserRepository)

	// モックの挙動を設定
	mockUserRepository.On("FetchUserById", "123").Return(nil, errors.New("failed to fetch user"))

	// サービス層メソッドの実行
	user, err := userService.FetchUserById("123")

	// エラーチェック
	assert.Error(t, err)

	// データが期待通りか確認
	assert.Nil(t, user)
	assert.Equal(t, "failed to fetch user", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockUserRepository.AssertExpectations(t)
}
