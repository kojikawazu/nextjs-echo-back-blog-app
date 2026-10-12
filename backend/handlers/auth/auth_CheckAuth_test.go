package handlers_auth

import (
	"backend/apperrors"
	"backend/config"
	"backend/models"
	services_auth "backend/services/auth"
	services_users "backend/services/blog_users"
	utils_cookie "backend/utils/cookie"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// generateTestToken はテスト用の認証トークンを生成するヘルパー
func generateTestToken(t *testing.T, expiresAt time.Time) string {
	t.Helper()
	tokenString, err := utils_cookie.NewAuthToken(&models.BlogUsersData{
		ID:    "user123",
		Email: "test@example.com",
		Name:  "Test User",
	}, expiresAt)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return tokenString
}

// generateLegacyToken は aud を持たない（#156 以前の形式の）トークンを生成するヘルパー
func generateLegacyToken(t *testing.T) string {
	t.Helper()
	claims := &models.Claims{
		UserID:           "user123",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.JwtKey)
	if err != nil {
		t.Fatalf("failed to generate legacy token: %v", err)
	}
	return tokenString
}

func TestHandler_CheckAuth(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(req *http.Request)
		expectedStatus int
		expectedBody   string
		wantErr        error // エラー系はハンドラの戻り値を検証する（ステータス・本文への変換は middlewares.HTTPErrorHandler）
	}{
		{
			name: "正常系_有効なJWTCookieで認証確認成功",
			setupCookie: func(req *http.Request) {
				tokenString := generateTestToken(t, time.Now().Add(1*time.Hour))
				req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"email":"test@example.com","message":"Authenticated","user_id":"user123","username":"Test User"}`,
		},
		{
			name:        "準正常系_Cookieなしの場合はErrAuthTokenMissingを返す",
			setupCookie: func(req *http.Request) {},
			wantErr:     apperrors.ErrAuthTokenMissing,
		},
		{
			name: "準正常系_不正なトークン値の場合はErrAuthTokenInvalidを返す",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "token", Value: "invalid-token-value"})
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			name: "準正常系_期限切れトークンの場合はErrAuthTokenInvalidを返す",
			setupCookie: func(req *http.Request) {
				tokenString := generateTestToken(t, time.Now().Add(-1*time.Hour))
				req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			// 誰でも発行できる訪問者トークンでは認証を通さない（#156。以前は 200 Authenticated を返していた）
			name: "準正常系_訪問者トークンの場合はErrAuthTokenInvalidを返す",
			setupCookie: func(req *http.Request) {
				tokenString, err := utils_cookie.NewVisitorToken("visitor-1", time.Now().Add(1*time.Hour))
				if err != nil {
					t.Fatalf("failed to generate visitor token: %v", err)
				}
				req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			// aud を持たない旧形式のトークンは受け付けない（#156 のデプロイで 1 回だけ再ログインが必要）
			name: "準正常系_audを持たない旧形式のトークンの場合はErrAuthTokenInvalidを返す",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "token", Value: generateLegacyToken(t)})
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/users/auth-check", nil)
			tt.setupCookie(req)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockAuthService := new(services_auth.MockAuthService)
			mockUserService := new(services_users.MockUserService)
			handler := NewAuthHandler(mockUserService, mockAuthService)

			err := handler.CheckAuth(c)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.False(t, c.Response().Committed)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
