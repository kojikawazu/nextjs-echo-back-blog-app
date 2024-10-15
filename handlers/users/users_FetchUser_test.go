package handlers_users

import (
	"backend/models"
	services_users "backend/services/users"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchUser(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/users/detail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUser := &models.UserData{
		ID:    "valid-user-id",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockUserService.On("FetchUserById", "valid-user-id").Return(mockUser, nil)

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	handler.FetchUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "John Doe")

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
}

func TestHandler_FetchUser_NotToken(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/users/detail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックの振る舞いを設定
	mockCookieUtils.On("GetAuthCookieValue", c).Return("", errors.New("no cookie"))

	// ハンドラーを実行
	err := handler.FetchUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error getting cookie")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertNotCalled(t, "FetchUserById", "valid-user-id")
}

func TestHandler_FetchUser_NotUserId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/users/detail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックの振る舞いを設定
	mockCookieUtils.On("GetAuthCookieValue", c).Return("mocked-token", nil)
	mockCookieUtils.On("GetUserIdFromToken", c, "mocked-token").Return("", errors.New("no user id"))

	// ハンドラーを実行
	err := handler.FetchUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error getting userId from token")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertNotCalled(t, "FetchUserById", "valid-user-id")
}

func TestHandler_FetchUser_NotUser(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/users/detail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("FetchUserById", "valid-user-id").Return(nil, errors.New("user not found"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err := handler.FetchUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "User not found")

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
}

func TestHandler_FetchUser_ServerError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/users/detail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("FetchUserById", "valid-user-id").Return(nil, errors.New("server error"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err := handler.FetchUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to fetch user")

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
}
