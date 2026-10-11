package handlers_auth

import (
	"backend/logger"
	"backend/models"
	services_auth "backend/services/auth"
	services_users "backend/services/blog_users"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	logTestEmail    = "secret-user@example.com"
	logTestPassword = "plain-secret-password"
)

// captureHandlerLogs は logger の INFO / ERROR 出力をバッファへ差し替え、テスト終了時に復元する。
func captureHandlerLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prevInfo, prevErr := logger.InfoLog.Writer(), logger.ErrorLog.Writer()
	logger.InfoLog.SetOutput(&buf)
	logger.ErrorLog.SetOutput(&buf)
	t.Cleanup(func() {
		logger.InfoLog.SetOutput(prevInfo)
		logger.ErrorLog.SetOutput(prevErr)
	})
	return &buf
}

func newLoginContext(body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

// TestHandler_Login_DoesNotLogCredentials はログイン時のログにメールアドレス・パスワードが残らないことを検証する。
func TestHandler_Login_DoesNotLogCredentials(t *testing.T) {
	body := `{"email":"` + logTestEmail + `","password":"` + logTestPassword + `"}`

	tests := []struct {
		name           string
		setupMock      func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService)
		expectedStatus int
		wantLog        string
	}{
		{
			// 正常系: 認証成功ログはユーザー ID で識別する
			name: "success_logs_user_id_only",
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", logTestEmail, logTestPassword).Return(nil)
				mockUser.On("FetchUserByEmailAndPassword", logTestEmail, logTestPassword).Return(&models.BlogUsersData{
					ID:       "user123",
					Email:    logTestEmail,
					Name:     "Test User",
					Password: logTestPassword,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			wantLog:        "User authenticated successfully: id=user123",
		},
		{
			// 準正常系: 認証失敗時のログにも資格情報が残らない
			name: "invalid_credentials",
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", logTestEmail, logTestPassword).Return(nil)
				mockUser.On("FetchUserByEmailAndPassword", logTestEmail, logTestPassword).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
			wantLog:        "Error fetching user: user not found",
		},
		{
			// 異常系: 想定外エラー時のログにも資格情報が残らない
			name: "unexpected_error",
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", logTestEmail, logTestPassword).Return(errors.New("unexpected error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := captureHandlerLogs(t)
			c, rec := newLoginContext(body)

			mockAuth := new(services_auth.MockAuthService)
			mockUser := new(services_users.MockUserService)
			tt.setupMock(mockAuth, mockUser)

			assert.NoError(t, NewAuthHandler(mockUser, mockAuth).Login(c))
			assert.Equal(t, tt.expectedStatus, rec.Code)

			logs := buf.String()
			assert.NotContains(t, logs, logTestEmail)
			assert.NotContains(t, logs, logTestPassword)
			if tt.wantLog != "" {
				assert.Contains(t, logs, tt.wantLog)
			}
		})
	}
}

// TestHandler_CheckAuth_DoesNotLogEmail は認証確認成功時のログがメールアドレスではなくユーザー ID を出すことを検証する（正常系）。
func TestHandler_CheckAuth_DoesNotLogEmail(t *testing.T) {
	buf := captureHandlerLogs(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/auth/check", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: generateTestToken(t, time.Now().Add(time.Hour))})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewAuthHandler(new(services_users.MockUserService), new(services_auth.MockAuthService))
	assert.NoError(t, handler.CheckAuth(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	logs := buf.String()
	// generateTestToken のメールアドレス
	assert.NotContains(t, logs, "test@example.com")
	assert.Contains(t, logs, "Authentication successful for user: id=user123")
}
