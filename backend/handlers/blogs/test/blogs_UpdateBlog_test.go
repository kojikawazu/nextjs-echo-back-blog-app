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

func TestHandler_UpdateBlog(t *testing.T) {
	// 正常なリクエストボディ
	validBody := map[string]string{"title": "Test Title", "githubUrl": "https://github.com", "category": "Tech", "description": "This is a test blog", "tags": "Go"}

	// withBody は正常なボディの 1 フィールドを差し替えたボディを返す
	withBody := func(key, value string) map[string]string {
		body := map[string]string{}
		for k, v := range validBody {
			body[k] = v
		}
		body[key] = value
		return body
	}

	// expectUpdate はボディの値で UpdateBlog が呼ばれるモックを設定する
	expectUpdate := func(m *service_blogs.MockBlogService, id string, body map[string]string, blog *models.BlogData, err error) {
		m.On("UpdateBlog", id, body["title"], body["githubUrl"], body["category"], body["description"], body["tags"]).Return(blog, err)
	}

	tests := []struct {
		name           string
		id             string
		requestBody    map[string]string // nil の場合は不正な JSON を送る
		setupMock      func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService)
		expectedStatus int    // 正常系のみ
		expectedBody   string // 正常系のみ（部分一致）
		wantErr        error  // 期待するエラー（nil は正常系）
	}{
		{
			name:        "正常系_全フィールド正常値でブログ更新成功",
			id:          "123",
			requestBody: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", validBody, &models.BlogData{Title: "Test Title"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `"title":"Test Title"`,
		},
		{
			name:        "準正常系_idが不正な場合ErrInvalidIDを返す",
			id:          "",
			requestBody: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "", validBody, nil, service_blogs.ErrInvalidID)
			},
			wantErr: service_blogs.ErrInvalidID,
		},
		{
			name:        "準正常系_titleが空文字の場合ErrInvalidTitleを返す",
			id:          "123",
			requestBody: withBody("title", ""),
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", withBody("title", ""), nil, service_blogs.ErrInvalidTitle)
			},
			wantErr: service_blogs.ErrInvalidTitle,
		},
		{
			name:        "準正常系_githubUrlが空文字の場合ErrInvalidGithubURLを返す",
			id:          "123",
			requestBody: withBody("githubUrl", ""),
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", withBody("githubUrl", ""), nil, service_blogs.ErrInvalidGithubURL)
			},
			wantErr: service_blogs.ErrInvalidGithubURL,
		},
		{
			name:        "準正常系_categoryが空文字の場合ErrInvalidCategoryを返す",
			id:          "123",
			requestBody: withBody("category", ""),
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", withBody("category", ""), nil, service_blogs.ErrInvalidCategory)
			},
			wantErr: service_blogs.ErrInvalidCategory,
		},
		{
			name:        "準正常系_descriptionが空文字の場合ErrInvalidDescriptionを返す",
			id:          "123",
			requestBody: withBody("description", ""),
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", withBody("description", ""), nil, service_blogs.ErrInvalidDescription)
			},
			wantErr: service_blogs.ErrInvalidDescription,
		},
		{
			name:        "準正常系_tagsが空文字の場合ErrInvalidTagsを返す",
			id:          "123",
			requestBody: withBody("tags", ""),
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", withBody("tags", ""), nil, service_blogs.ErrInvalidTags)
			},
			wantErr: service_blogs.ErrInvalidTags,
		},
		{
			name:        "準正常系_ブログが存在しない場合ErrBlogNotFoundを返す",
			id:          "123",
			requestBody: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", validBody, nil, service_blogs.ErrBlogNotFound)
			},
			wantErr: service_blogs.ErrBlogNotFound,
		},
		{
			name:        "準正常系_認証Cookie取得に失敗した場合ErrAuthTokenMissingを返す",
			id:          "123",
			requestBody: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				mockCookie.On("GetAuthCookieValue", c, "token").Return("", errors.New("no cookie"))
			},
			wantErr: apperrors.ErrAuthTokenMissing,
		},
		{
			name:        "準正常系_JWT解析に失敗した場合ErrAuthTokenInvalidを返す",
			id:          "123",
			requestBody: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				mockCookie.On("GetAuthCookieValue", c, "token").Return("invalid-token", nil)
				mockCookie.On("GetUserIdFromToken", c, "invalid-token").Return("", errors.New("invalid token"))
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			name:        "異常系_不正なJSONボディの場合ErrInvalidRequestBodyを返す",
			id:          "123",
			requestBody: nil,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
			},
			wantErr: apperrors.ErrInvalidRequestBody,
		},
		{
			name:        "異常系_Serviceが予期しないエラーを返す場合そのエラーを返す",
			id:          "123",
			requestBody: validBody,
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				expectUpdate(mockService, "123", validBody, nil, errServiceFailure)
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

			req := httptest.NewRequest(http.MethodPut, "/api/blogs/update/"+tt.id, body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockBlogService := new(service_blogs.MockBlogService)
			handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils, mockBlogService)

			err := handler.UpdateBlog(c)

			if tt.wantErr != nil {
				assertHandlerError(t, err, rec, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}

			mockCookieUtils.AssertExpectations(t)
			mockBlogService.AssertExpectations(t)
		})
	}
}
