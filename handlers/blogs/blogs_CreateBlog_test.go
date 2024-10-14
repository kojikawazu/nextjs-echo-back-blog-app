package handlers_blogs

import (
	"backend/models"
	service_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateBlog(t *testing.T) {
	e := echo.New()

	// フォームデータを設定
	form := url.Values{}
	form.Add("title", "Test Title")
	form.Add("githubUrl", "https://github.com")
	form.Add("category", "Tech")
	form.Add("description", "This is a test blog")
	form.Add("tags", "Go")

	// JWT の署名キーを設定し、正しいトークンを生成
	token := "mocked-token"
	validUserId := "valid-user-id"

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/blogs/create", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	cookie := &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// サービスとハンドラーをモックする
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockBlogService := new(service_blogs.MockBlogService)
	handler := NewBlogHandler(mockBlogService, mockCookieUtils)

	// モックの振る舞いを設定
	mockCookieUtils.On("GetAuthCookieValue", c).Return(token, nil)
	mockCookieUtils.On("GetUserIdFromToken", c, token).Return(validUserId, nil)
	mockBlogService.On("CreateBlog", validUserId, "Test Title", "https://github.com", "Tech", "This is a test blog", "Go").Return(models.BlogData{
		Title: "Test Title",
	}, nil)

	// テストを実行
	err := handler.CreateBlog(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test Title")

	// モックの呼び出しを確認
	mockCookieUtils.AssertExpectations(t)
	mockBlogService.AssertExpectations(t)
}
