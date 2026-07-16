package handlers_blog_likes

import (
	services_blogs_likes "backend/services/blog_likes"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_DeleteBlogLike(t *testing.T) {
	tests := []struct {
		name           string
		blogId         string
		setupMock      func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "正常系_いいね削除成功",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("DeleteBlogLike", "blog-1", "valid-visit-id").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Blog like deleted successfully"}`,
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
			name:   "異常系_Serviceがエラーを返す場合500を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("DeleteBlogLike", "blog-1", "valid-visit-id").Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Error deleting blog like"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/blog-likes/delete/"+tt.blogId, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("blogId")
			c.SetParamValues(tt.blogId)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(services_blogs_likes.MockBlogLikeService)
			handler := NewBlogLikeHandler(mockService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils, mockService)

			err := handler.DeleteBlogLike(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockCookieUtils.AssertExpectations(t)
			mockService.AssertExpectations(t)
		})
	}
}
