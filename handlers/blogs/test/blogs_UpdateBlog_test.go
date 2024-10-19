package handlers_blogs_test

import (
	handlers_blogs "backend/handlers/blogs"
	"backend/models"
	service_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateBlog(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "https://github.com",
		"category":    "Tech",
		"description": "This is a test blog",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "Test Title", "https://github.com", "Tech", "This is a test blog", "Go").Return(&models.BlogData{
		Title: "Test Title",
	}, nil)

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test Title")

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_InvalidId(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "https://github.com",
		"category":    "Tech",
		"description": "This is a test blog",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "", "Test Title", "https://github.com", "Tech", "This is a test blog", "Go").Return(nil, errors.New("invalid id"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid id"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_InvalidTitle(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "",
		"githubUrl":   "https://github.com",
		"category":    "Tech",
		"description": "This is a test blog",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "", "https://github.com", "Tech", "This is a test blog", "Go").Return(nil, errors.New("invalid title"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid title"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_InvalidGithubUrl(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "",
		"category":    "Tech",
		"description": "This is a test blog",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "Test Title", "", "Tech", "This is a test blog", "Go").Return(nil, errors.New("invalid githubUrl"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid githubUrl"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_InvalidCategory(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "https://github.com",
		"category":    "",
		"description": "This is a test blog",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "Test Title", "https://github.com", "", "This is a test blog", "Go").Return(nil, errors.New("invalid category"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid category"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_InvalidDescription(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "https://github.com",
		"category":    "Tech",
		"description": "",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "Test Title", "https://github.com", "Tech", "", "Go").Return(nil, errors.New("invalid description"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid description"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_InvalidTags(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "https://github.com",
		"category":    "Tech",
		"description": "This is a test blog",
		"tags":        "",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "Test Title", "https://github.com", "Tech", "This is a test blog", "").Return(nil, errors.New("invalid tags"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid tags"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_NoUpdate(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "https://github.com",
		"category":    "Tech",
		"description": "This is a test blog",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "Test Title", "https://github.com", "Tech", "This is a test blog", "Go").Return(nil, errors.New("failed to update blog"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error":"Failed to update blog"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}

func TestHandler_UpdateBlog_ServerError(t *testing.T) {
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"title":       "Test Title",
		"githubUrl":   "https://github.com",
		"category":    "Tech",
		"description": "This is a test blog",
		"tags":        "Go",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/update/123", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("id")
	c.SetParamValues("123")

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("UpdateBlog", "123", "Test Title", "https://github.com", "Tech", "This is a test blog", "Go").Return(nil, errors.New("server error"))

	// モッククッキーを設定
	handlers_blogs.SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.UpdateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error":"Server error"}`, rec.Body.String())

	// モックの呼び出しを確認
	mockBlogService.AssertExpectations(t)
}
