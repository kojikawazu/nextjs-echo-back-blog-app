package handlers_auth

import (
	"backend/config"
	"backend/models"
	services_auth "backend/services/auth"
	services_users "backend/services/blog_users"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// generateTestToken はテスト用JWTトークンを生成するヘルパー
func generateTestToken(t *testing.T, expiresAt time.Time) string {
	t.Helper()
	claims := &models.Claims{
		UserID:   "user123",
		Email:    "test@example.com",
		Username: "Test User",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return tokenString
}

func TestHandler_CheckAuth(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(req *http.Request)
		expectedStatus int
		expectedBody   string
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
			name:           "準正常系_Cookieなしの場合401を返す",
			setupCookie:    func(req *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Token not found"}`,
		},
		{
			name: "準正常系_不正なトークン値の場合401を返す",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "token", Value: "invalid-token-value"})
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Invalid token"}`,
		},
		{
			name: "準正常系_期限切れトークンの場合401を返す",
			setupCookie: func(req *http.Request) {
				tokenString := generateTestToken(t, time.Now().Add(-1*time.Hour))
				req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"message":"Invalid token"}`,
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

			handler.CheckAuth(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
