package handlers_auth

import (
	services_auth "backend/services/auth"
	services_users "backend/services/blog_users"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Logout(t *testing.T) {
	tests := []struct {
		name           string
		setupCookie    func(req *http.Request)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常系_ログアウト成功しCookieが削除される",
			setupCookie: func(req *http.Request) {
				req.AddCookie(&http.Cookie{Name: "token", Value: "some-token"})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Logout successful"}`,
		},
		{
			name:           "準正常系_Cookieなしでもログアウト成功",
			setupCookie:    func(req *http.Request) {},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Logout successful"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/users/logout", nil)
			tt.setupCookie(req)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockAuthService := new(services_auth.MockAuthService)
			mockUserService := new(services_users.MockUserService)
			handler := NewAuthHandler(mockUserService, mockAuthService)

			handler.Logout(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			// tokenクッキーが削除されているか確認（Expiresが過去日時）
			for _, cookie := range rec.Result().Cookies() {
				if cookie.Name == "token" {
					assert.True(t, cookie.Expires.Before(time.Now()),
						"token cookie should have a past expiration time")
				}
			}
		})
	}
}
