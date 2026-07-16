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

func TestHandler_FetchBlogLikesByVisitId(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "正常系_訪問IDに紐づくいいね一覧を取得成功",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("FetchBlogLikesByVisitId", "valid-visit-id").Return([]models.BlogLikesData{
					{ID: "1", BlogId: "1", VisitId: "valid-visit-id"},
					{ID: "2", BlogId: "2", VisitId: "valid-visit-id"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":"1","blog_id":"1","visit_id":"valid-visit-id","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},{"id":"2","blog_id":"2","visit_id":"valid-visit-id","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "正常系_いいねが0件の場合空配列を返す",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("FetchBlogLikesByVisitId", "valid-visit-id").Return([]models.BlogLikesData{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
		{
			name: "準正常系_visit-id-token Cookieがない場合500を返す",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return("", errors.New("no cookie"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to get visit id token"}`,
		},
		{
			name: "準正常系_visitId取得に失敗した場合500を返す",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				token := "mocked-token"
				req.AddCookie(&http.Cookie{Name: "visit-id-token", Value: token})
				mockCookie.On("GetAuthCookieValue", c, "visit-id-token").Return(token, nil)
				mockCookie.On("GetVisitIdFromToken", c, token).Return("", errors.New("no visit id"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to get visit id"}`,
		},
		{
			name: "異常系_Serviceがエラーを返す場合500を返す",
			setupMock: func(c echo.Context, req *http.Request, mockCookie *utils_cookie.MockCookieUtils, mockService *services_blogs_likes.MockBlogLikeService) {
				SetMockBlogCookies(c, req, mockCookie)
				mockService.On("FetchBlogLikesByVisitId", "valid-visit-id").Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Error fetching blog likes by visit id"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/blog-likes", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(services_blogs_likes.MockBlogLikeService)
			handler := NewBlogLikeHandler(mockService, mockCookieUtils)

			tt.setupMock(c, req, mockCookieUtils, mockService)

			err := handler.FetchBlogLikesByVisitId(c)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())

			mockCookieUtils.AssertExpectations(t)
			mockService.AssertExpectations(t)
		})
	}
}
