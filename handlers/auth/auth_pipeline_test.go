package handlers_auth

import (
	"backend/models"
	services_auth "backend/services/auth"
	services_users "backend/services/blog_users"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestHandler_AuthPipeline は Login → CheckAuth → Logout の一連フローを検証する統合テスト
func TestHandler_AuthPipeline(t *testing.T) {
	SetupTest(t)

	mockUser := &models.BlogUsersData{
		ID:        "user123",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("正常系_Login後にCheckAuthが成功しLogoutで認証が無効になる", func(t *testing.T) {
		e := echo.New()
		mockAuthService := new(services_auth.MockAuthService)
		mockUserService := new(services_users.MockUserService)
		handler := NewAuthHandler(mockUserService, mockAuthService)

		// --- Step 1: Login ---
		mockAuthService.On("Login", "test@example.com", "password123").Return(nil)
		mockUserService.On("FetchUserByEmailAndPassword", "test@example.com", "password123").Return(mockUser, nil)

		loginBody := `{"email":"test@example.com","password":"password123"}`
		loginReq := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(loginBody))
		loginReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		loginRec := httptest.NewRecorder()
		loginCtx := e.NewContext(loginReq, loginRec)

		err := handler.Login(loginCtx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, loginRec.Code)
		assert.JSONEq(t, `{"message":"Login successful"}`, loginRec.Body.String())

		// Loginレスポンスからtokenクッキーを取得
		var tokenCookie *http.Cookie
		for _, c := range loginRec.Result().Cookies() {
			if c.Name == "token" {
				tokenCookie = c
				break
			}
		}
		assert.NotNil(t, tokenCookie, "token cookie should be set after login")

		// --- Step 2: CheckAuth (有効なトークン) ---
		checkReq := httptest.NewRequest(http.MethodGet, "/api/users/auth-check", nil)
		checkReq.AddCookie(tokenCookie)
		checkRec := httptest.NewRecorder()
		checkCtx := e.NewContext(checkReq, checkRec)

		handler.CheckAuth(checkCtx)
		assert.Equal(t, http.StatusOK, checkRec.Code)
		assert.JSONEq(t, `{"message":"Authenticated","user_id":"user123","username":"Test User","email":"test@example.com"}`, checkRec.Body.String())

		// --- Step 3: Logout ---
		logoutReq := httptest.NewRequest(http.MethodPost, "/api/users/logout", nil)
		logoutRec := httptest.NewRecorder()
		logoutCtx := e.NewContext(logoutReq, logoutRec)

		err = handler.Logout(logoutCtx)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, logoutRec.Code)
		assert.JSONEq(t, `{"message":"Logout successful"}`, logoutRec.Body.String())

		// ログアウト後のCookieがMaxAge=-1（削除済み）であることを確認
		var logoutCookie *http.Cookie
		for _, c := range logoutRec.Result().Cookies() {
			if c.Name == "token" {
				logoutCookie = c
				break
			}
		}
		if logoutCookie != nil {
			assert.True(t, logoutCookie.MaxAge < 0 || logoutCookie.Value == "", "token cookie should be invalidated after logout")
		}

		mockAuthService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})
}
