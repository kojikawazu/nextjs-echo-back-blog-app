package handlers_blogs_test

import (
	handlers_blogs "backend/handlers/blogs"
	"backend/models"
	service_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchBlogsByUserId(t *testing.T) {
	e := echo.New()
	// パスパラメータとして userId を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/user/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("userId")
	c.SetParamValues("1")

	// モックサービスをインスタンス化
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockService := new(service_blogs.MockBlogService)
	handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

	// モックデータの設定
	mockBlogs := []models.BlogData{
		{
			ID:         "1",
			BlogUserId: "1",
			Title:      "title1",
			GithubUrl:  "https://github.com/user/repo1",
			Category:   "Category1",
			Tags:       "Tag1",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "2",
			BlogUserId: "1",
			Title:      "title2",
			GithubUrl:  "https://github.com/user/repo2",
			Category:   "Category2",
			Tags:       "Tag2",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
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

	// サービス層のメソッドが呼ばれたことを確認
	mockService.AssertExpectations(t)
}

// TestHandler_FetchBlogsByUserId_Errors はサービス層・パスパラメータの検証で失敗した場合に、ハンドラがエラーを返す（レスポンスは書き込まない）ことを検証する。
func TestHandler_FetchBlogsByUserId_Errors(t *testing.T) {
	tests := []struct {
		name      string
		param     string
		setupMock func(mockService *service_blogs.MockBlogService)
		wantErr   error
	}{
		{
			name:  "準正常系_userIdが不正な場合ErrInvalidUserIDを返す",
			param: "",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogsByUserId", "").Return(nil, service_blogs.ErrInvalidUserID)
			},
			wantErr: service_blogs.ErrInvalidUserID,
		},
		{
			name:  "異常系_Serviceが予期しないエラーを返す場合そのエラーを返す",
			param: "1",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogsByUserId", "1").Return(nil, errServiceFailure)
			},
			wantErr: errServiceFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/blogs/user/"+tt.param, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("userId")
			c.SetParamValues(tt.param)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(service_blogs.MockBlogService)
			handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)
			tt.setupMock(mockService)

			err := handler.FetchBlogsByUserId(c)

			assertHandlerError(t, err, rec, tt.wantErr)
			mockService.AssertExpectations(t)
		})
	}
}
