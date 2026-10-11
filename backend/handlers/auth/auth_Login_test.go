package handlers_auth

import (
	"backend/apperrors"
	"backend/models"
	services_auth "backend/services/auth"
	services_users "backend/services/blog_users"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// errUnexpected はサービス層の想定外エラー（DB 障害等）を表すテスト用エラー。
var errUnexpected = errors.New("unexpected error")

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		contentType    string
		setupMock      func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService)
		expectedStatus int
		expectedBody   string
		wantErr        error // エラー系はハンドラの戻り値を検証する（ステータス・本文への変換は middlewares.HTTPErrorHandler）
		wantNotErr     error // 戻り値が一致してはならない sentinel（nil のときは検証しない）
	}{
		{
			name:        "正常系_正しいメールとパスワードでログイン成功",
			requestBody: map[string]string{"email": "test@example.com", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "test@example.com", "password123").Return(nil)
				mockUser.On("FetchUserByEmailAndPassword", "test@example.com", "password123").Return(&models.BlogUsersData{
					ID:    "user123",
					Email: "test@example.com",
					Name:  "Test User",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Login successful"}`,
		},
		{
			name:        "準正常系_メールが空文字の場合はErrEmailAndPasswordRequiredを返す",
			requestBody: map[string]string{"email": "", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "", "password123").Return(services_auth.ErrEmailAndPasswordRequired)
			},
			wantErr: services_auth.ErrEmailAndPasswordRequired,
		},
		{
			name:        "準正常系_パスワードが空文字の場合はErrEmailAndPasswordRequiredを返す",
			requestBody: map[string]string{"email": "test@example.com", "password": ""},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "test@example.com", "").Return(services_auth.ErrEmailAndPasswordRequired)
			},
			wantErr: services_auth.ErrEmailAndPasswordRequired,
		},
		{
			name:        "準正常系_メールフォーマットが不正の場合はErrInvalidEmailFormatを返す",
			requestBody: map[string]string{"email": "not-an-email", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "not-an-email", "password123").Return(services_auth.ErrInvalidEmailFormat)
			},
			wantErr: services_auth.ErrInvalidEmailFormat,
		},
		{
			name:        "準正常系_認証失敗時はErrInvalidCredentialsを返す（ユーザー列挙対策）",
			requestBody: map[string]string{"email": "notexist@example.com", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "notexist@example.com", "password123").Return(nil)
				mockUser.On("FetchUserByEmailAndPassword", "notexist@example.com", "password123").Return(nil, services_users.ErrInvalidCredentials)
			},
			wantErr: services_users.ErrInvalidCredentials,
		},
		{
			name:        "異常系_UserServiceの想定外エラーは認証失敗に丸めずそのまま返す",
			requestBody: map[string]string{"email": "test@example.com", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "test@example.com", "password123").Return(nil)
				mockUser.On("FetchUserByEmailAndPassword", "test@example.com", "password123").Return(nil, errUnexpected)
			},
			wantErr:    errUnexpected,
			wantNotErr: services_users.ErrInvalidCredentials,
		},
		{
			name:        "異常系_AuthServiceが予期しないエラーを返す場合はそのまま返す",
			requestBody: map[string]string{"email": "test@example.com", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "test@example.com", "password123").Return(errUnexpected)
			},
			wantErr: errUnexpected,
		},
		{
			name:        "異常系_不正なJSONボディの場合はErrInvalidRequestBodyを返す",
			requestBody: nil,
			contentType: echo.MIMEApplicationJSON,
			setupMock:   func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {},
			wantErr:     apperrors.ErrInvalidRequestBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			var body *bytes.Buffer
			if tt.requestBody == nil {
				// 不正なJSONを送信
				body = bytes.NewBufferString("{invalid json}")
			} else {
				jsonData, err := json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				body = bytes.NewBuffer(jsonData)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/users/login", body)
			req.Header.Set(echo.HeaderContentType, tt.contentType)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockAuthService := new(services_auth.MockAuthService)
			mockUserService := new(services_users.MockUserService)
			tt.setupMock(mockAuthService, mockUserService)

			handler := NewAuthHandler(mockUserService, mockAuthService)
			err := handler.Login(c)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				if tt.wantNotErr != nil {
					assert.NotErrorIs(t, err, tt.wantNotErr)
				}
				// エラー時はハンドラ自身でレスポンスを書かない
				assert.False(t, c.Response().Committed)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				var resp map[string]string
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Login successful", resp["message"])
			}

			mockAuthService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}
