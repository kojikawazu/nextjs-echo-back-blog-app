package handlers_blogs_test

import (
	"backend/apperrors"
	handlers_blogs "backend/handlers/blogs"
	service_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_DeleteBlog(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		setupMock func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService)
		wantErr   error // 期待するエラー（nil は正常系）
	}{
		{
			name: "正常系_ブログ削除成功で204を返す",
			id:   "123",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("DeleteBlog", "123").Return(nil)
			},
		},
		{
			name: "準正常系_idが不正な場合ErrInvalidIDを返す",
			id:   "",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("DeleteBlog", "").Return(service_blogs.ErrInvalidID)
			},
			wantErr: service_blogs.ErrInvalidID,
		},
		{
			name: "準正常系_認証Cookie取得に失敗した場合ErrAuthTokenMissingを返す",
			id:   "123",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				mockCookie.On("GetAuthCookieValue", c, "token").Return("", errors.New("no cookie"))
			},
			wantErr: apperrors.ErrAuthTokenMissing,
		},
		{
			name: "準正常系_JWT解析に失敗した場合ErrAuthTokenInvalidを返す",
			id:   "123",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				mockCookie.On("GetAuthCookieValue", c, "token").Return("invalid-token", nil)
				mockCookie.On("GetUserIdFromToken", c, "invalid-token").Return("", errors.New("invalid token"))
			},
			wantErr: apperrors.ErrAuthTokenInvalid,
		},
		{
			name: "異常系_Serviceが予期しないエラーを返す場合そのエラーを返す",
			id:   "123",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *service_blogs.MockBlogService) {
				handlers_blogs.SetMockBlogCookies(c, req, mockCookie)
				mockService.On("DeleteBlog", "123").Return(errServiceFailure)
			},
			wantErr: errServiceFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/api/blogs/delete/"+tt.id, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockBlogService := new(service_blogs.MockBlogService)
			handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils, mockBlogService)

			err := handler.DeleteBlog(c)

			if tt.wantErr != nil {
				assertHandlerError(t, err, rec, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusNoContent, rec.Code)
				assert.Empty(t, rec.Body.String())
			}

			mockCookieUtils.AssertExpectations(t)
			mockBlogService.AssertExpectations(t)
		})
	}
}
