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
	"github.com/stretchr/testify/mock"
)

func TestHandler_FetchBlogPopular(t *testing.T) {
	e := echo.New()

	// パスパラメータとして count を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/popular/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("count")
	c.SetParamValues("2")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockBlogs := []models.BlogData{
		{
			ID:        "1",
			UserId:    "1",
			Title:     "title1",
			GithubUrl: "https://github.com/user/repo1",
			Category:  "Category1",
			Tags:      "Tag1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			UserId:    "1",
			Title:     "title2",
			GithubUrl: "https://github.com/user/repo2",
			Category:  "Category2",
			Tags:      "Tag2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.On("FetchBlogPopular", 2).Return(mockBlogs, nil)

	// ハンドラーを実行
	err := handler.FetchBlogPopular(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "title1")
	assert.Contains(t, rec.Body.String(), "title2")

	// サービス層のメソッドが呼ばれたことを確認
	mockService.AssertExpectations(t)
}

// userIdが空の場合の異常系
func TestHandler_FetchBlogPopular_EmptyCount(t *testing.T) {
	e := echo.New()

	// パスパラメータが空のリクエストを作成
	req := httptest.NewRequest(http.MethodGet, "/api/blog/popular/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを空に設定
	c.SetParamNames("count")
	c.SetParamValues("")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックサービスの設定（ブログが見つからない場合）
	mockService.On("FetchBlogPopular", 0).Return(nil, errors.New("invalid count"))

	// ハンドラーを実行
	err := handler.FetchBlogPopular(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid count")

	// モックの呼び出しがないことを確認
	mockService.AssertNotCalled(t, "FetchBlogPopular", mock.Anything)
}

// ブログが見つからない場合の異常系
func TestHandler_FetchBlogPopular_BlogNotFound(t *testing.T) {
	e := echo.New()

	// パスパラメータとして count を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/popular/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// パスパラメータを設定
	c.SetParamNames("count")
	c.SetParamValues("1")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックサービスの設定（ブログが見つからない場合）
	mockService.On("FetchBlogPopular", 1).Return(nil, errors.New("blog not found"))

	// ハンドラーを実行
	err := handler.FetchBlogPopular(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "blog not found")

	// モックの呼び出しを検証
	mockService.AssertExpectations(t)
}

// サービス層でエラーが発生した場合の異常系
func TestHandler_FetchBlogPopular_ServiceError(t *testing.T) {
	e := echo.New()

	// パスパラメータとして count を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/popular/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("count")
	c.SetParamValues("1")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックサービスの設定（一般的なエラーが発生した場合）
	mockService.On("FetchBlogPopular", 1).Return(nil, errors.New("Error fetching popular blogs"))

	// ハンドラーを実行
	err := handler.FetchBlogPopular(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error fetching popular blogs")

	// モックの呼び出しを検証
	mockService.AssertExpectations(t)
}
