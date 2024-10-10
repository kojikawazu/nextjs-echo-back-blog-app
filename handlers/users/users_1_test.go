package handlers_users

import (
	services_users "backend/services/users"

	"backend/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserByEmailAndPassword(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// モックデータの設定
	mockUser := &models.UserData{
		ID:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockService.On("FetchUserByEmailAndPassword", "john@example.com", "password123").Return(mockUser, nil)

	// ハンドラーを実行
	handler.GetUserByEmailAndPassword(c)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "John Doe")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetUserByEmailAndPassword_ValidationError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"email":"", "password":""}` // 空のemailとpassword
	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// モックの挙動を設定（バリデーションエラーを返す）
	mockService.On("FetchUserByEmailAndPassword", "", "").Return(nil, errors.New("email and password are required"))

	// ハンドラーを実行
	handler.GetUserByEmailAndPassword(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Email and password are required")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetUserByEmailAndPassword_InvalidEmailFormat(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"email":"invalid-email", "password":"password123"}` // 無効なメールフォーマット
	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// モックの挙動を設定（無効なメールフォーマットの場合のエラーを返す）
	mockService.On("FetchUserByEmailAndPassword", "invalid-email", "password123").Return(nil, errors.New("invalid email format"))

	// ハンドラーを実行
	handler.GetUserByEmailAndPassword(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid email format")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetUserByEmailAndPassword_UserNotFound(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// サービスがユーザーが見つからないことを返すようにモックの挙動を設定
	mockService.On("FetchUserByEmailAndPassword", "john@example.com", "password123").Return(nil, errors.New("user not found"))

	// ハンドラーを実行
	handler.GetUserByEmailAndPassword(c)

	// ステータスコードとレスポンス内容を確認
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "User not found")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_GetUserByEmailAndPassword_ServiceError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	body := `{"email":"john@example.com", "password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/user", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(services_users.MockUserService)
	handler := NewUserHandler(mockService)

	// サービスがエラーを返すようにモックの挙動を設定
	mockService.On("FetchUserByEmailAndPassword", "john@example.com", "password123").Return(nil, errors.New("error fetching user"))

	// ハンドラーを実行
	handler.GetUserByEmailAndPassword(c)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to fetch user")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}
