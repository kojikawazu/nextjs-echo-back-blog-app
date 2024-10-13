package handlers_blogs

import (
	"backend/models"
	service_blogs "backend/services/blogs"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchBlogsByUserId(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/blog/user", strings.NewReader(`{"user_id": "1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// モックデータの設定
	mockBlogs := []models.BlogData{
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
			UserId:    "1",
			Title:     "title2",
			GithubUrl: "https://github.com/user/repo2",
			Category:  "Category2",
			Tag:       "Tag2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockService.On("FetchBlogsByUserId", "1").Return(mockBlogs, nil)

	// ハンドラーを実行
	err := handler.FetchBlogsByUserId(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "title1")
	assert.Contains(t, rec.Body.String(), "title2")

	// サービス層のメソッドが呼ばれていないことを確認
	mockService.AssertExpectations(t)
}

// リクエストボディが無効な場合の異常系
func TestHandler_FetchBlogsByUserId_InvalidRequestBody(t *testing.T) {
	e := echo.New()
	// 無効なボディを渡す
	req := httptest.NewRequest(http.MethodPost, "/api/blog/user", strings.NewReader(`invalid_json_body`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// ハンドラーを実行
	err := handler.FetchBlogsByUserId(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid request body")

	// サービス層のメソッドが呼ばれていないことを確認
	mockService.AssertNotCalled(t, "FetchBlogsByUserId")
}

// userIdが空の場合の異常系
func TestHandler_FetchBlogsByUserId_EmptyUserId(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/blog/user", strings.NewReader(`{"user_id": ""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// モックサービスの設定（ブログが見つからない場合）
	mockService.On("FetchBlogsByUserId", "").Return(nil, errors.New("invalid userId"))

	// ハンドラーを実行
	err := handler.FetchBlogsByUserId(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid userId")

	// サービス層のメソッドが呼ばれていないことを確認
	mockService.AssertExpectations(t)
}

// ブログが見つからない場合の異常系
func TestHandler_FetchBlogsByUserId_BlogNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/blog/user", strings.NewReader(`{"user_id": "1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// モックサービスの設定（ブログが見つからない場合）
	mockService.On("FetchBlogsByUserId", "1").Return(nil, errors.New("blog not found"))

	// ハンドラーを実行
	err := handler.FetchBlogsByUserId(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Blog not found")

	mockService.AssertExpectations(t)
}

// サービス層でエラーが発生した場合の異常系
func TestHandler_FetchBlogsByUserId_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/blog/user", strings.NewReader(`{"user_id": "1"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockService)

	// モックサービスの設定（一般的なエラーが発生した場合）
	mockService.On("FetchBlogsByUserId", "1").Return(nil, errors.New("some internal error"))

	// ハンドラーを実行
	err := handler.FetchBlogsByUserId(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error fetching blog")

	mockService.AssertExpectations(t)
}
