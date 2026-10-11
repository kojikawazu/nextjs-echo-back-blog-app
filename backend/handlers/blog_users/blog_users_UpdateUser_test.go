package handlers_blog_users

import (
	"backend/apperrors"
	"backend/models"
	services_users "backend/services/blog_users"
	utils_cookie "backend/utils/cookie"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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
	handler := NewBlogUsersHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUser := &models.BlogUsersData{
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
	err = handler.UpdateBlogUsers(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "John Doe")

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
	mockCookieUtils.AssertExpectations(t)
}

// TestHandler_UpdateUser_Errors はユーザー更新の失敗時にハンドラが返すエラーを検証する。
// ステータス・本文への変換は middlewares.HTTPErrorHandler の責務のため、ここでは戻り値の sentinel を検証する。
func TestHandler_UpdateUser_Errors(t *testing.T) {
	errUnexpected := errors.New("server error")
	validBody := `{"name":"John Doe","email":"john@example.com","password":"password","newPassword":"new-password"}`
	longPassword := strings.Repeat("a", 73)
	mockUser := &models.BlogUsersData{ID: "valid-user-id", Name: "John Doe", Email: "john@example.com"}

	// serviceReturns は認証クッキーが有効で、UpdateUser が err を返すモック設定を作る。
	serviceReturns := func(name, email, password, newPassword string, err error) func(echo.Context, *http.Request, *utils_cookie.MockCookieUtils, *services_users.MockUserService) {
		return func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
			SetMockBlogCookies(c, req, mockCookieUtils)
			mockUserService.On("UpdateUser", "valid-user-id", name, email, password, newPassword).Return(nil, err)
		}
	}
	bodyOf := func(name, email, password, newPassword string) string {
		b, err := json.Marshal(map[string]string{"name": name, "email": email, "password": password, "newPassword": newPassword})
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}
		return string(b)
	}

	tests := []struct {
		name          string
		body          string
		setupMock     func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService)
		wantErr       error
		expectService bool
	}{
		{
			// 準正常系: 認証クッキーが無い → ErrAuthTokenMissing
			name: "token_missing",
			body: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				mockCookieUtils.On("GetAuthCookieValue", c, "token").Return("", errors.New("no cookie"))
			},
			wantErr: apperrors.ErrAuthTokenMissing,
		},
		{
			// 準正常系: トークンからユーザー ID を取得できない → ErrAuthTokenInvalid
			name: "token_invalid",
			body: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				mockCookieUtils.On("GetAuthCookieValue", c, "token").Return("mocked-token", nil)
				mockCookieUtils.On("GetUserIdFromToken", c, "mocked-token").Return("", errors.New("no user id"))
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			// 準正常系: JSON として解釈できないボディ → ErrInvalidRequestBody
			name: "invalid_body",
			body: "{invalid json}",
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				SetMockBlogCookies(c, req, mockCookieUtils)
			},
			wantErr: apperrors.ErrInvalidRequestBody,
		},
		{
			// 準正常系: name が空
			name:          "name_required",
			body:          bodyOf("", "john@example.com", "password", "new-password"),
			setupMock:     serviceReturns("", "john@example.com", "password", "new-password", services_users.ErrNameRequired),
			wantErr:       services_users.ErrNameRequired,
			expectService: true,
		},
		{
			// 準正常系: email が空
			name:          "email_required",
			body:          bodyOf("John Doe", "", "password", "new-password"),
			setupMock:     serviceReturns("John Doe", "", "password", "new-password", services_users.ErrEmailRequired),
			wantErr:       services_users.ErrEmailRequired,
			expectService: true,
		},
		{
			// 準正常系: email の形式が不正
			name:          "invalid_email_format",
			body:          bodyOf("John Doe", "test", "password", "new-password"),
			setupMock:     serviceReturns("John Doe", "test", "password", "new-password", services_users.ErrInvalidEmailFormat),
			wantErr:       services_users.ErrInvalidEmailFormat,
			expectService: true,
		},
		{
			// 準正常系: 現在のパスワードが空
			name:          "password_required",
			body:          bodyOf("John Doe", "john@example.com", "", "new-password"),
			setupMock:     serviceReturns("John Doe", "john@example.com", "", "new-password", services_users.ErrPasswordRequired),
			wantErr:       services_users.ErrPasswordRequired,
			expectService: true,
		},
		{
			// 準正常系: 新しいパスワードが空
			name:          "new_password_required",
			body:          bodyOf("John Doe", "john@example.com", "password", ""),
			setupMock:     serviceReturns("John Doe", "john@example.com", "password", "", services_users.ErrNewPasswordRequired),
			wantErr:       services_users.ErrNewPasswordRequired,
			expectService: true,
		},
		{
			// 準正常系: 新しいパスワードが 72 バイトを超える
			name:          "new_password_too_long",
			body:          bodyOf("John Doe", "john@example.com", "password", longPassword),
			setupMock:     serviceReturns("John Doe", "john@example.com", "password", longPassword, services_users.ErrPasswordTooLong),
			wantErr:       services_users.ErrPasswordTooLong,
			expectService: true,
		},
		{
			// 準正常系: 現在のパスワード不一致
			name:          "invalid_current_password",
			body:          validBody,
			setupMock:     serviceReturns("John Doe", "john@example.com", "password", "new-password", services_users.ErrInvalidCurrentPassword),
			wantErr:       services_users.ErrInvalidCurrentPassword,
			expectService: true,
		},
		{
			// 準正常系: ユーザー不在
			name:          "user_not_found",
			body:          validBody,
			setupMock:     serviceReturns("John Doe", "john@example.com", "password", "new-password", services_users.ErrUserNotFound),
			wantErr:       services_users.ErrUserNotFound,
			expectService: true,
		},
		{
			// 異常系: サービスの想定外エラーをそのまま返す（500 への変換は HTTPErrorHandler）
			name:          "unexpected_error",
			body:          validBody,
			setupMock:     serviceReturns("John Doe", "john@example.com", "password", "new-password", errUnexpected),
			wantErr:       errUnexpected,
			expectService: true,
		},
		{
			// 異常系: 更新後のトークン再作成に失敗した場合は原因を辿れる形で返す
			name: "create_token_error",
			body: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				SetMockBlogCookies(c, req, mockCookieUtils)
				mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "john@example.com", "password", "new-password").Return(mockUser, nil)
				mockCookieUtils.On("GetAuthCookieExpirationTime").Return(time.Now().Add(1 * time.Hour))
				mockCookieUtils.On("CreateToken", mockUser).Return("", errUnexpected)
			},
			wantErr:       errUnexpected,
			expectService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/api/users/update", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockUserService := new(services_users.MockUserService)
			tt.setupMock(c, req, mockCookieUtils, mockUserService)

			err := NewBlogUsersHandler(mockUserService, mockCookieUtils).UpdateBlogUsers(c)

			assert.ErrorIs(t, err, tt.wantErr)
			// エラー時はハンドラ自身でレスポンスを書かず、クッキーも更新しない
			assert.False(t, c.Response().Committed)
			mockCookieUtils.AssertNotCalled(t, "UpdateAuthCookie", mock.Anything, mock.Anything, mock.Anything)
			mockCookieUtils.AssertExpectations(t)
			if tt.expectService {
				mockUserService.AssertExpectations(t)
			} else {
				mockUserService.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}
