package handlers_blog_likes

import (
	services_blogs_likes "backend/services/blog_likes"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GenerateVisitorId(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常系_VisitIDが存在しない場合新規生成してCookieをセット",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils) {
				// Cookieなし → GetAuthCookieValue がエラー
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return("", errors.New("no cookie"))
				expirationTime := time.Now().Add(1 * time.Hour)
				mockCookie.On("GetAuthCookieExpirationTime").Return(expirationTime)
				mockCookie.On("CreateVisitIdToken").Return("new-visit-token", nil)
				mockCookie.On("AddVisitIdCoookie", c, "new-visit-token", expirationTime).Return()
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Visitor id generated successfully"}`,
		},
		{
			name: "正常系_VisitIDがすでに存在する場合スキップ",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils) {
				token := "existing-token"
				req.AddCookie(&http.Cookie{Name: "visit-id-token", Value: token})
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return(token, nil)
				mockCookie.On("GetVisitIdFromToken", c, token).Return("existing-visit-id", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Visitor id already exists"}`,
		},
		{
			name: "異常系_トークン生成に失敗した場合500を返す",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils) {
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return("", errors.New("no cookie"))
				expirationTime := time.Now().Add(1 * time.Hour)
				mockCookie.On("GetAuthCookieExpirationTime").Return(expirationTime)
				mockCookie.On("CreateVisitIdToken").Return("", errors.New("token generation failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to create visitor token"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/blog-likes/generate-visit-id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(services_blogs_likes.MockBlogLikeService)
			handler := NewBlogLikeHandler(mockService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils)

			err := handler.GenerateVisitorId(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockCookieUtils.AssertExpectations(t)
		})
	}
}
