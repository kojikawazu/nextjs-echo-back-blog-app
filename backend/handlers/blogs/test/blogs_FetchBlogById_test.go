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
		ID:         "1",
		BlogUserId: "1",
		Title:      "title1",
		GithubUrl:  "",
		Category:   "Category1",
		Tags:       "Tag1",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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

// TestHandler_FetchBlogById_Errors はサービス層・パスパラメータの検証で失敗した場合に、ハンドラがエラーを返す（レスポンスは書き込まない）ことを検証する。
func TestHandler_FetchBlogById_Errors(t *testing.T) {
	tests := []struct {
		name      string
		param     string
		setupMock func(mockService *service_blogs.MockBlogService)
		wantErr   error
	}{
		{
			name:  "準正常系_idが不正な場合ErrInvalidIDを返す",
			param: "",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogById", "").Return(nil, service_blogs.ErrInvalidID)
			},
			wantErr: service_blogs.ErrInvalidID,
		},
		{
			name:  "準正常系_ブログが存在しない場合ErrBlogNotFoundを返す",
			param: "1",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogById", "1").Return(nil, service_blogs.ErrBlogNotFound)
			},
			wantErr: service_blogs.ErrBlogNotFound,
		},
		{
			name:  "異常系_Serviceが予期しないエラーを返す場合そのエラーを返す",
			param: "1",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogById", "1").Return(nil, errServiceFailure)
			},
			wantErr: errServiceFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/blogs/detail/"+tt.param, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.param)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(service_blogs.MockBlogService)
			handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)
			tt.setupMock(mockService)

			err := handler.FetchBlogById(c)

			assertHandlerError(t, err, rec, tt.wantErr)
			mockService.AssertExpectations(t)
		})
	}
}
