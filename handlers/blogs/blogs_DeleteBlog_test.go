package handlers_blogs

import (
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
	e := echo.New()

	// リクエストを作成
	req := httptest.NewRequest(http.MethodDelete, "/blogs/delete/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("DeleteBlog", "123").Return(nil, nil)

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err := handler.DeleteBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Empty(t, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_DeleteBlog_InvalidId(t *testing.T) {
	e := echo.New()

	// リクエストを作成
	req := httptest.NewRequest(http.MethodDelete, "/blogs/delete/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("DeleteBlog", "").Return(errors.New("invalid id"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err := handler.DeleteBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid id"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_DeleteBlog_NoDelete(t *testing.T) {
	e := echo.New()

	// リクエストを作成
	req := httptest.NewRequest(http.MethodDelete, "/blogs/delete/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("DeleteBlog", "123").Return(errors.New("failed to delete blog"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err := handler.DeleteBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error":"Failed to delete blog"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_DeleteBlog_ServerError(t *testing.T) {
	e := echo.New()

	// リクエストを作成
	req := httptest.NewRequest(http.MethodDelete, "/blogs/delete/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("DeleteBlog", "123").Return(errors.New("server error"))

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err := handler.DeleteBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error":"Server error"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}
