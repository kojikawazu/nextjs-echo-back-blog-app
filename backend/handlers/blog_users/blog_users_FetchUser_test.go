package handlers_blog_users

import (
	"backend/apperrors"
	"backend/models"
	services_users "backend/services/blog_users"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	handler := NewBlogUsersHandler(mockUserService, mockCookieUtils)

	// モックデータの設定
	mockUser := &models.BlogUsersData{
		ID:    "valid-user-id",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockUserService.On("FetchUserById", "valid-user-id").Return(mockUser, nil)

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err := handler.FetchBlogUsers(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "John Doe")

	// モックが期待通りに呼び出されたかを確認
	mockUserService.AssertExpectations(t)
}

// TestHandler_FetchUser_Errors はユーザー取得の失敗時にハンドラが返すエラーを検証する。
// ステータス・本文への変換は middlewares.HTTPErrorHandler の責務のため、ここでは戻り値の sentinel を検証する。
func TestHandler_FetchUser_Errors(t *testing.T) {
	errUnexpected := errors.New("server error")

	tests := []struct {
		name          string
		setupMock     func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService)
		wantErr       error
		expectService bool
	}{
		{
			// 準正常系: 認証クッキーが無い → ErrAuthTokenMissing
			name: "token_missing",
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				mockCookieUtils.On("GetAuthCookieValue", c, "token").Return("", errors.New("no cookie"))
			},
			wantErr: apperrors.ErrAuthTokenMissing,
		},
		{
			// 準正常系: トークンからユーザー ID を取得できない → ErrAuthTokenInvalid
			name: "token_invalid",
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				mockCookieUtils.On("GetAuthCookieValue", c, "token").Return("mocked-token", nil)
				mockCookieUtils.On("GetUserIdFromToken", c, "mocked-token").Return("", errors.New("no user id"))
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			// 準正常系: ユーザー不在 → サービスの ErrUserNotFound をそのまま返す
			name: "user_not_found",
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				SetMockBlogCookies(c, req, mockCookieUtils)
				mockUserService.On("FetchUserById", "valid-user-id").Return(nil, services_users.ErrUserNotFound)
			},
			wantErr:       services_users.ErrUserNotFound,
			expectService: true,
		},
		{
			// 異常系: サービスの想定外エラーをそのまま返す（500 への変換は HTTPErrorHandler）
			name: "unexpected_error",
			setupMock: func(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils, mockUserService *services_users.MockUserService) {
				SetMockBlogCookies(c, req, mockCookieUtils)
				mockUserService.On("FetchUserById", "valid-user-id").Return(nil, errUnexpected)
			},
			wantErr:       errUnexpected,
			expectService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/users/detail", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockUserService := new(services_users.MockUserService)
			tt.setupMock(c, req, mockCookieUtils, mockUserService)

			err := NewBlogUsersHandler(mockUserService, mockCookieUtils).FetchBlogUsers(c)

			assert.ErrorIs(t, err, tt.wantErr)
			// エラー時はハンドラ自身でレスポンスを書かない
			assert.False(t, c.Response().Committed)
			mockCookieUtils.AssertExpectations(t)
			if tt.expectService {
				mockUserService.AssertExpectations(t)
			} else {
				mockUserService.AssertNotCalled(t, "FetchUserById", mock.Anything)
			}
		})
	}
}
