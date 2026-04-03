package handlers_blog_likes

import (
	"backend/models"
	services_blogs_likes "backend/services/blog_likes"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateBlogLike(t *testing.T) {
	tests := []struct {
		name           string
		blogId         string
		setupMock      func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "正常系_いいね作成成功",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlogLike", "blog-1", "valid-visit-id").Return(&models.BlogLikesData{
					ID:      "like-1",
					BlogId:  "blog-1",
					VisitId: "valid-visit-id",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"like-1","blog_id":"blog-1","visit_id":"valid-visit-id","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:   "準正常系_BlogIdまたはVisitIdが空の場合400を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlogLike", "blog-1", "valid-visit-id").Return(nil, errors.New("BlogId or VisitId is empty"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"BlogId or VisitId is empty"}`,
		},
		{
			name:   "準正常系_すでにいいね済みの場合400を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlogLike", "blog-1", "valid-visit-id").Return(nil, errors.New("blog is already liked"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Blog is already liked"}`,
		},
		{
			name:   "準正常系_visit-id-token Cookieがない場合500を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return("", errors.New("no cookie"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to get visit id token"}`,
		},
		{
			name:   "準正常系_visitId取得に失敗した場合500を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				token := "mocked-token"
				req.AddCookie(&http.Cookie{Name: "visit-id-token", Value: token})
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return(token, nil)
				mockCookie.On("GetVisitIdFromToken", c, token).Return("", errors.New("parse error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to get visit id"}`,
		},
		{
			name:   "異常系_Serviceが予期しないエラーを返す場合500を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("CreateBlogLike", "blog-1", "valid-visit-id").Return(nil, errors.New("db connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Error creating blog like"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/blog-likes/create/"+tt.blogId, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("blogId")
			c.SetParamValues(tt.blogId)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(services_blogs_likes.MockBlogLikeService)
			handler := NewBlogLikeHandler(mockService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils, mockService)

			err := handler.CreateBlogLike(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockCookieUtils.AssertExpectations(t)
			mockService.AssertExpectations(t)
		})
	}
}
