package handlers_blogs_test

import (
	handlers_blogs "backend/handlers/blogs"
	"backend/models"
	service_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchBlogs(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		setupMock      func(mockService *service_blogs.MockBlogService)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
	}{
		{
			name: "正常系_ブログ一覧を複数件取得成功",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogs").Return([]models.BlogData{
					{ID: "1", BlogUserId: "u1", Title: "title1", GithubUrl: "https://github.com/user/repo1", Category: "Category1", Tags: "Tag1", CreatedAt: now, UpdatedAt: now},
					{ID: "2", BlogUserId: "u2", Title: "title2", GithubUrl: "https://github.com/user/repo2", Category: "Category2", Tags: "Tag2", CreatedAt: now, UpdatedAt: now},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, `"title":"title1"`)
				assert.Contains(t, body, `"title":"title2"`)
			},
		},
		{
			name: "正常系_ブログが0件の場合空配列を返す",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogs").Return([]models.BlogData{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `[]`, body)
			},
		},
		{
			name: "異常系_Serviceがエラーを返す場合500を返す",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogs").Return(nil, errors.New("db connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"error":"Error fetching blogs"}`, body)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/blogs", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(service_blogs.MockBlogService)
			handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

			tt.setupMock(mockService)

			err := handler.FetchBlogs(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.checkBody(t, rec.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
