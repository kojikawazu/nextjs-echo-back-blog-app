package handlers_blog_likes

import (
	services_blogs_likes "backend/services/blog_likes"
	utils_cookie "backend/utils/cookie"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_IsBlogLiked(t *testing.T) {
	tests := []struct {
		name           string
		blogId         string
		setupMock      func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService)
		expectedStatus int
		expectedBody   string
		wantErr        error // 非 nil のときはハンドラがこのエラーを返すことを検証する
	}{
		{
			name:   "正常系_いいね済みの場合isLiked:trueを返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("IsBlogLiked", "blog-1", "valid-visit-id").Return(true, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"isLiked":true}`,
		},
		{
			name:   "正常系_いいね未済の場合isLiked:falseを返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("IsBlogLiked", "blog-1", "valid-visit-id").Return(false, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"isLiked":false}`,
		},
		{
			name:   "準正常系_blogIdが空の場合400を返す",
			blogId: "",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("IsBlogLiked", "", "valid-visit-id").Return(false, services_blogs_likes.ErrBlogIDOrVisitIDEmpty)
			},
			wantErr: services_blogs_likes.ErrBlogIDOrVisitIDEmpty,
		},
		{
			// 確認の失敗を「いいねしていない」として 200 で隠さない
			name:   "異常系_Serviceがエラーを返す場合500を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("IsBlogLiked", "blog-1", "valid-visit-id").Return(false, errService)
			},
			wantErr: errService,
		},
		{
			name:   "準正常系_visit-id-token Cookieがない場合500を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return("", errNoCookie)
			},
			wantErr: errNoCookie,
		},
		{
			name:   "準正常系_visitId取得に失敗した場合500を返す",
			blogId: "blog-1",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				token := "mocked-token"
				req.AddCookie(&http.Cookie{Name: "visit-id-token", Value: token})
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return(token, nil)
				mockCookie.On("GetVisitIdFromToken", c, token).Return("", errParseToken)
			},
			wantErr: errParseToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/blog-likes/is-liked/"+tt.blogId, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("blogId")
			c.SetParamValues(tt.blogId)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(services_blogs_likes.MockBlogLikeService)
			handler := NewBlogLikeHandler(mockService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils, mockService)

			err := handler.IsBlogLiked(c)
			if tt.wantErr != nil {
				assertHandlerError(t, err, rec, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}

			mockCookieUtils.AssertExpectations(t)
			mockService.AssertExpectations(t)
		})
	}
}
