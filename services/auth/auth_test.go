package services_auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_Login_Success(t *testing.T) {
	// モックの作成
	mockAuthService := new(MockAuthService)

	// 期待する入力と出力の設定
	email := "test@example.com"
	password := "password123"
	mockAuthService.On("Login", email, password).Return(nil)

	// テスト対象の関数呼び出し
	err := mockAuthService.Login(email, password)

	// アサーション
	assert.NoError(t, err)
	mockAuthService.AssertExpectations(t)
}

func TestService_Login_Failure(t *testing.T) {
	// モックの作成
	mockAuthService := new(MockAuthService)

	// 期待する入力と出力の設定
	email := "test@example.com"
	password := "wrongpassword"
	mockAuthService.On("Login", email, password).Return(errors.New("invalid credentials"))

	// テスト対象の関数呼び出し
	err := mockAuthService.Login(email, password)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	mockAuthService.AssertExpectations(t)
}
