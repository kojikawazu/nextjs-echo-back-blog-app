package handlers_blogs

import (
	"backend/models"
	service_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateBlog(t *testing.T) {
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

	// JWT の署名キーを設定し、正しいトークンを生成
	validUserId := "valid-user-id"

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/create", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockBlogService.On("CreateBlog", validUserId, "Test Title", "https://github.com", "Tech", "This is a test blog", "Go").Return(models.BlogData{
		Title: "Test Title",
	}, nil)

	// モッククッキーを設定
	SetMockBlogCookies(c, req, mockCookieUtils)

	// テストを実行
	err = handler.CreateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test Title")

	// モックの呼び出しを確認
	mockCookieUtils.AssertExpectations(t)
	mockBlogService.AssertExpectations(t)
}
