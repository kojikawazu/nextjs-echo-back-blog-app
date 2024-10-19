package handlers_users

import (
	"backend/models"
	services_users "backend/services/users"
	utils_cookie "backend/utils/cookie"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_UpdateUser(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "john@example.com",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "john@example.com", "password", "new-password").Return(mockUser, nil)

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// CookieUtilsのモック設定(追加分)
	mockCookieUtils.On("GetAuthCookieExpirationTime").Return(time.Now().Add(1 * time.Hour))
	mockCookieUtils.On("CreateToken", mockUser).Return("new-mocked-token", nil)
	mockCookieUtils.On("UpdateAuthCookie", c, "new-mocked-token", mock.Anything).Return(nil)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "John Doe")

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
	mockCookieUtils.AssertExpectations(t)
}

func TestHandler_UpdateUser_NotToken(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "john@example.com",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックの振る舞いを設定
	mockCookieUtils.On("GetAuthCookieValue", c, "token").Return("", errors.New("no cookie"))

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error getting cookie")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertNotCalled(t, "UpdateUser", "valid-user-id", "John Doe", "john@example.com", "password", "new-password")
}

func TestHandler_UpdateUser_NotUserId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "john@example.com",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックの振る舞いを設定
	mockCookieUtils.On("GetAuthCookieValue", c, "token").Return("mocked-token", nil)
	mockCookieUtils.On("GetUserIdFromToken", c, "mocked-token").Return("", errors.New("no user id"))

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error getting userId from token")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertNotCalled(t, "UpdateUser", "valid-user-id", "John Doe", "john@example.com", "password", "new-password")
}

func TestHandler_UpdateUser_InvalidName(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "",
		"email":       "john@example.com",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("UpdateUser", "valid-user-id", "", "john@example.com", "password", "new-password").Return(nil, errors.New("name is required"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Name is required")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestHandler_UpdateUser_InvalidEmail01(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "", "password", "new-password").Return(nil, errors.New("email is required"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email is required")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestHandler_UpdateUser_InvalidEmail02(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "test",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "test", "password", "new-password").Return(nil, errors.New("invalid email format"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid email format")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestHandler_UpdateUser_InvalidPassword(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "john@example.com",
		"password":    "",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "john@example.com", "", "new-password").Return(nil, errors.New("password is required"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Password is required")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestHandler_UpdateUser_InvalidNewPassword(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "john@example.com",
		"password":    "password",
		"newPassword": "",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "john@example.com", "password", "").Return(nil, errors.New("new password is required"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "New password is required")

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
}

func TestHandler_UpdateUser_InvalidUser(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "john@example.com",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "john@example.com", "password", "new-password").Return(nil, errors.New("user not found"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "User not found")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

func TestHandler_UpdateUser_NotUpdated(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"name":        "John Doe",
		"email":       "john@example.com",
		"password":    "password",
		"newPassword": "new-password",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	handler := NewUserHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "john@example.com", "password", "new-password").Return(nil, errors.New("server error"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err = handler.UpdateUser(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to update user")

	// モックが期待通りに呼び出されたかを確認
	mockCookieUtils.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}
