package handlers_auth

import (
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

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		contentType    string
		setupMock      func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService)
		expectedStatus int
		expectedBody   string
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
			name:        "準正常系_メールが空文字の場合400を返す",
			requestBody: map[string]string{"email": "", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "", "password123").Return(errors.New("email and password are required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Email and password are required"}`,
		},
		{
			name:        "準正常系_パスワードが空文字の場合400を返す",
			requestBody: map[string]string{"email": "test@example.com", "password": ""},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "test@example.com", "").Return(errors.New("email and password are required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Email and password are required"}`,
		},
		{
			name:        "準正常系_メールフォーマットが不正の場合400を返す",
			requestBody: map[string]string{"email": "not-an-email", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "not-an-email", "password123").Return(errors.New("invalid email format"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid email format"}`,
		},
		{
			name:        "準正常系_認証失敗時に401を返す（ユーザー列挙対策）",
			requestBody: map[string]string{"email": "notexist@example.com", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "notexist@example.com", "password123").Return(nil)
				mockUser.On("FetchUserByEmailAndPassword", "notexist@example.com", "password123").Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid credentials"}`,
		},
		{
			name:        "異常系_AuthServiceが予期しないエラーを返す場合500を返す",
			requestBody: map[string]string{"email": "test@example.com", "password": "password123"},
			contentType: echo.MIMEApplicationJSON,
			setupMock: func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {
				mockAuth.On("Login", "test@example.com", "password123").Return(errors.New("unexpected error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"An error occurred"}`,
		},
		{
			name:           "異常系_不正なJSONボディの場合400を返す",
			requestBody:    nil,
			contentType:    echo.MIMEApplicationJSON,
			setupMock:      func(mockAuth *services_auth.MockAuthService, mockUser *services_users.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
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
			handler.Login(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus == http.StatusOK {
				var resp map[string]string
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Login successful", resp["message"])
			} else {
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}

			mockAuthService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}
