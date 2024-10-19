package handlers_blogs_likes

import (
	"backend/models"
	services_blogs_likes "backend/services/blogs_likes"
	utils_cookie "backend/utils/cookie"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchBlogLikesByVisitId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blog-likes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(services_blogs_likes.MockBlogLikeService)
	handler := NewBlogLikeHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockBlog := []models.BlogLikeData{
		{
			ID:      "1",
			BlogId:  "1",
			VisitId: "valid-visit-id",
		},
		{
			ID:      "2",
			BlogId:  "2",
			VisitId: "valid-visit-id",
		},
	}
	mockService.On("FetchBlogLikesByVisitId", "valid-visit-id").Return(mockBlog, nil)

	// クッキーのモックを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err := handler.FetchBlogLikesByVisitId(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "valid-visit-id")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_FetchBlogLikesByVisitId_NoCookie(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blog-likes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(services_blogs_likes.MockBlogLikeService)
	handler := NewBlogLikeHandler(mockService, mockCookieUtils)

	// モックを設定
	mockCookieUtils.On("GetAuthCookieValue", c, "visit-id-token").Return("", errors.New("no cookie"))

	// ハンドラーを実行
	err := handler.FetchBlogLikesByVisitId(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to get visit id token")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertNotCalled(t, "FetchBlogLikesByVisitId", "valid-visit-id")
}

func TestHandler_FetchBlogLikesByVisitId_NoVisitId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blog-likes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// リクエストにクッキーを追加
	token := "mocked-token"
	cookie := &http.Cookie{
		Name:  "visit-id-token",
		Value: token,
		Path:  "/",
	}
	req.AddCookie(cookie)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(services_blogs_likes.MockBlogLikeService)
	handler := NewBlogLikeHandler(mockService, mockCookieUtils)

	// モックを設定
	mockCookieUtils.On("GetAuthCookieValue", c, "visit-id-token").Return("mocked-token", nil)
	mockCookieUtils.On("GetVisitIdFromToken", c, token).Return("", errors.New("no visit id"))

	// ハンドラーを実行
	err := handler.FetchBlogLikesByVisitId(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to get visit id")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertNotCalled(t, "FetchBlogLikesByVisitId", "valid-visit-id")
}

func TestHandler_FetchBlogLikesByVisitId_NotData(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blog-likes", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(services_blogs_likes.MockBlogLikeService)
	handler := NewBlogLikeHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockService.On("FetchBlogLikesByVisitId", "valid-visit-id").Return(nil, errors.New("no data"))

	// クッキーのモックを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// ハンドラーを実行
	err := handler.FetchBlogLikesByVisitId(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error fetching blog likes by visit id")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}
