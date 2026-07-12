package handlers_blogs_test

import (
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

func TestHandler_FetchTags(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blogs/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockTags := []string{"Tag1", "Tag2", "Tag3"}
	mockService.On("FetchBlogTags").Return(mockTags, nil)

	// ハンドラーを実行
	err := handler.FetchBlogTags(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `["Tag1","Tag2","Tag3"]`, rec.Body.String())

	// モックの呼び出しを検証
	mockService.AssertExpectations(t)
}

func TestHandler_FetchTags_NoData(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blogs/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockTags := []string{}
	mockService.On("FetchBlogTags").Return(mockTags, nil)

	// ハンドラーを実行
	err := handler.FetchBlogTags(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())

	// モックの呼び出しを検証
	mockService.AssertExpectations(t)
}

func TestHandler_FetchTags_ErrorCase(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/blogs/tags", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockService.On("FetchBlogTags").Return(nil, errors.New("some error occurred"))

	// ハンドラーを実行
	err := handler.FetchBlogTags(c)

	// ステータスコードとレスポンス内容の確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error fetching tags")

	// モックの呼び出しを検証
	mockService.AssertExpectations(t)
}
