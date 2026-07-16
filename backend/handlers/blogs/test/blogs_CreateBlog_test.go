package handlers_blogs_test

import (
	handlers_blogs "backend/handlers/blogs"
	"backend/models"
	service_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateBlog(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "正常系_全フィールド正常値でブログ作成成功",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "This is a test blog", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "Tech", "This is a test blog", "Go").Return(&models.BlogData{
					Title: "Test Title", GithubUrl: "https://github.com", Category: "Tech", Description: "This is a test blog", Tags: "Go",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"","blog_user_id":"","title":"Test Title","github_url":"https://github.com","category":"Tech","description":"This is a test blog","tags":"Go","likes":0,"comment_cnt":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:        "準正常系_titleが空文字の場合400を返す",
			requestBody: map[string]string{"title": "", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "", "https://github.com", "Tech", "test", "Go").Return(nil, errors.New("invalid title"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid title"}`,
		},
		{
			name:        "準正常系_githubUrlが空文字の場合400を返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "", "Tech", "test", "Go").Return(nil, errors.New("invalid githubUrl"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid githubUrl"}`,
		},
		{
			name:        "準正常系_categoryが空文字の場合400を返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "", "test", "Go").Return(nil, errors.New("invalid category"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid category"}`,
		},
		{
			name:        "準正常系_descriptionが空文字の場合400を返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "Tech", "", "Go").Return(nil, errors.New("invalid description"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid description"}`,
		},
		{
			name:        "準正常系_tagsが空文字の場合400を返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": ""},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "Tech", "test", "").Return(nil, errors.New("invalid tags"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid tags"}`,
		},
		{
			name:        "準正常系_認証Cookie取得に失敗した場合401を返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				mockCookie.On("GetAuthCookieValue", c, "token").Return("", errors.New("no cookie"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Error getting cookie"}`,
		},
		{
			name:        "準正常系_JWT解析に失敗した場合401を返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				token := "invalid-token"
				req.AddCookie(&http.Cookie{Name: "token", Value: token})
				mockCookie.On("GetAuthCookieValue", c, "token").Return(token, nil)
				mockCookie.On("GetUserIdFromToken", c, token).Return("", errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Error getting userId from token"}`,
		},
		{
			name:        "異常系_不正なJSONボディの場合400を返す",
			requestBody: nil,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid request body"}`,
		},
		{
			name:        "異常系_Serviceが予期しないエラーを返す場合500を返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "Tech", "test", "Go").Return(nil, errors.New("db connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			var body *bytes.Buffer
			if tt.requestBody == nil {
				body = bytes.NewBufferString("{invalid json}")
			} else {
				jsonData, err := json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				body = bytes.NewBuffer(jsonData)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/blogs/create", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockBlogService := new(service_blogs.MockBlogService)
			handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils, mockBlogService)

			err := handler.CreateBlog(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockCookieUtils.AssertExpectations(t)
			mockBlogService.AssertExpectations(t)
		})
	}
}
