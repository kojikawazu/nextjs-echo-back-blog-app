package handlers_blogs_test

import (
	"backend/apperrors"
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
		expectedStatus int    // 正常系のみ
		expectedBody   string // 正常系のみ
		wantErr        error  // 期待するエラー（nil は正常系）
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
			name:        "準正常系_titleが空文字の場合ErrInvalidTitleを返す",
			requestBody: map[string]string{"title": "", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "", "https://github.com", "Tech", "test", "Go").Return(nil, service_blogs.ErrInvalidTitle)
			},
			wantErr: service_blogs.ErrInvalidTitle,
		},
		{
			name:        "準正常系_githubUrlが空文字の場合ErrInvalidGithubURLを返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "", "Tech", "test", "Go").Return(nil, service_blogs.ErrInvalidGithubURL)
			},
			wantErr: service_blogs.ErrInvalidGithubURL,
		},
		{
			name:        "準正常系_categoryが空文字の場合ErrInvalidCategoryを返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "", "test", "Go").Return(nil, service_blogs.ErrInvalidCategory)
			},
			wantErr: service_blogs.ErrInvalidCategory,
		},
		{
			name:        "準正常系_descriptionが空文字の場合ErrInvalidDescriptionを返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "Tech", "", "Go").Return(nil, service_blogs.ErrInvalidDescription)
			},
			wantErr: service_blogs.ErrInvalidDescription,
		},
		{
			name:        "準正常系_tagsが空文字の場合ErrInvalidTagsを返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": ""},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "Tech", "test", "").Return(nil, service_blogs.ErrInvalidTags)
			},
			wantErr: service_blogs.ErrInvalidTags,
		},
		{
			name:        "準正常系_認証Cookie取得に失敗した場合ErrAuthTokenMissingを返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				mockCookie.On("GetAuthCookieValue", c, "token").Return("", errors.New("no cookie"))
			},
			wantErr: apperrors.ErrAuthTokenMissing,
		},
		{
			name:        "準正常系_JWT解析に失敗した場合ErrAuthTokenInvalidを返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				token := "invalid-token"
				req.AddCookie(&http.Cookie{Name: "token", Value: token})
				mockCookie.On("GetAuthCookieValue", c, "token").Return(token, nil)
				mockCookie.On("GetUserIdFromToken", c, token).Return("", errors.New("invalid token"))
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			name:        "異常系_不正なJSONボディの場合ErrInvalidRequestBodyを返す",
			requestBody: nil,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
			},
			wantErr: apperrors.ErrInvalidRequestBody,
		},
		{
			name:        "異常系_Serviceが予期しないエラーを返す場合そのエラーを返す",
			requestBody: map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "test", "tags": "Go"},
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlog", "valid-user-id", "Test Title", "https://github.com", "Tech", "test", "Go").Return(nil, errServiceFailure)
			},
			wantErr: errServiceFailure,
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

			if tt.wantErr != nil {
				assertHandlerError(t, err, rec, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}

			mockCookieUtils.AssertExpectations(t)
			mockBlogService.AssertExpectations(t)
		})
	}
}
