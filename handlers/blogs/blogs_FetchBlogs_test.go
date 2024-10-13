package handlers_blogs

import (
	service_blogs "backend/services/blogs"
	"errors"
	"time"

	"backend/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchBlogs(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// モックデータの設定
	mockBlog := []models.BlogData{
		{
			ID:        "1",
			UserId:    "1",
			Title:     "title1",
			GithubUrl: "https://github.com/user/repo1",
			Category:  "Category1",
			Tag:       "Tag1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			UserId:    "2",
			Title:     "title2",
			GithubUrl: "https://github.com/user/repo2",
			Category:  "Category2",
			Tag:       "Tag2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.On("FetchBlogs").Return(mockBlog, nil)

	// ハンドラーを実行
	err := handler.FetchBlogs(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "title2")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_FetchBlogs_Error(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// サービス層がエラーを返すように設定
	mockService.On("FetchBlogs").Return(nil, errors.New("some error occurred"))

	// ハンドラーを実行
	err := handler.FetchBlogs(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error fetching blogs")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}

func TestHandler_FetchBlogs_Empty(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blogs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// サービス層が空のブログリストを返すように設定
	mockService.On("FetchBlogs").Return([]models.BlogData{}, nil)

	// ハンドラーを実行
	err := handler.FetchBlogs(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "[]")

	// モックが期待通りに呼び出されたかを確認
	mockService.AssertExpectations(t)
}
