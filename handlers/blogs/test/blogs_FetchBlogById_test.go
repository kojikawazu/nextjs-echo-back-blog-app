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

func TestHandler_FetchBlogById(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	// パスパラメータとして userId を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/detail/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("1")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockBlog := &models.BlogData{
		ID:        "1",
		UserId:    "1",
		Title:     "title1",
		GithubUrl: "",
		Category:  "Category1",
		Tags:      "Tag1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockService.On("FetchBlogById", "1").Return(mockBlog, nil)

	// ハンドラーを実行
	err := handler.FetchBlogById(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "title1")

	// サービス層のメソッドが呼ばれたことを確認
	mockService.AssertExpectations(t)
}

func TestHandler_FetchBlogById_InvalidId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	// パスパラメータとして userId を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/detail/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockService.On("FetchBlogById", "").Return(nil, errors.New("invalid id"))

	// ハンドラーを実行
	err := handler.FetchBlogById(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid id")

	// サービス層のメソッドが呼ばれたことを確認
	mockService.AssertExpectations(t)
}

func TestHandler_FetchBlogById_InvalidNotBlog(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	// パスパラメータとして userId を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/detail/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("1")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定

	mockService.On("FetchBlogById", "1").Return(nil, errors.New("blog not found"))

	// ハンドラーを実行
	err := handler.FetchBlogById(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Blog not found")

	// サービス層のメソッドが呼ばれたことを確認
	mockService.AssertExpectations(t)
}
