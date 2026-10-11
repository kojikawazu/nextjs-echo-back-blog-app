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

func TestHandler_FetchBlogPopular(t *testing.T) {
	e := echo.New()

	// パスパラメータとして count を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/blog/popular/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("count")
	c.SetParamValues("2")

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
	mockService.On("FetchBlogPopular", 2).Return(mockBlogs, nil)

	// ハンドラーを実行
	err := handler.FetchBlogPopular(c)

	// ステータスコードとレスポンスの確認
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "title1")
	assert.Contains(t, rec.Body.String(), "title2")

	// サービス層のメソッドが呼ばれたことを確認
	mockService.AssertExpectations(t)
}

// TestHandler_FetchBlogPopular_Errors はサービス層・パスパラメータの検証で失敗した場合に、ハンドラがエラーを返す（レスポンスは書き込まない）ことを検証する。
func TestHandler_FetchBlogPopular_Errors(t *testing.T) {
	tests := []struct {
		name      string
		param     string
		setupMock func(mockService *service_blogs.MockBlogService)
		wantErr   error
	}{
		{
			name:  "準正常系_countが空文字の場合サービスを呼ばずErrInvalidCountを返す",
			param: "",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				// 数値変換に失敗するためサービスは呼ばれない
			},
			wantErr: service_blogs.ErrInvalidCount,
		},
		{
			name:  "準正常系_countが数値でない場合サービスを呼ばずErrInvalidCountを返す",
			param: "abc",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				// 数値変換に失敗するためサービスは呼ばれない
			},
			wantErr: service_blogs.ErrInvalidCount,
		},
		{
			name:  "準正常系_countが0以下の場合ErrInvalidCountを返す",
			param: "0",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogPopular", 0).Return(nil, service_blogs.ErrInvalidCount)
			},
			wantErr: service_blogs.ErrInvalidCount,
		},
		{
			name:  "異常系_Serviceが予期しないエラーを返す場合そのエラーを返す",
			param: "1",
			setupMock: func(mockService *service_blogs.MockBlogService) {
				mockService.On("FetchBlogPopular", 1).Return(nil, errServiceFailure)
			},
			wantErr: errServiceFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/blogs/popular/"+tt.param, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("count")
			c.SetParamValues(tt.param)

			mockCookieUtils := new(utils_cookie.MockCookieUtils)
			mockService := new(service_blogs.MockBlogService)
			handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)
			tt.setupMock(mockService)

			err := handler.FetchBlogPopular(c)

			assertHandlerError(t, err, rec, tt.wantErr)
			mockService.AssertExpectations(t)
		})
	}
}
