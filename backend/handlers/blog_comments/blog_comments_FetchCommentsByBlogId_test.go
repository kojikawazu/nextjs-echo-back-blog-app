package handlers_blog_comments

import (
	"backend/models"
	services_comments "backend/services/blog_comments"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchCommentsByBlogId(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		blogId         string
		setupMock      func(mockService *services_comments.MockCommentService)
		expectedStatus int
		checkBody      func(t *testing.T, body string)
		wantErr        error // 非 nil のときはハンドラがこのエラーを返すことを検証する
	}{
		{
			name:   "正常系_ブログIDに紐づくコメント一覧を取得成功",
			blogId: "1",
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("FetchCommentsByBlogId", "1").Return([]models.BlogCommentsData{
					{ID: "1", BlogId: "1", GuestUser: "guestUser1", Comment: "comment1", CreatedAt: now},
					{ID: "2", BlogId: "1", GuestUser: "guestUser2", Comment: "comment2", CreatedAt: now},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.Contains(t, body, "guestUser1")
				assert.Contains(t, body, "guestUser2")
				assert.Contains(t, body, "comment1")
			},
		},
		{
			name:   "正常系_コメントが0件の場合空配列を返す",
			blogId: "1",
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("FetchCommentsByBlogId", "1").Return([]models.BlogCommentsData{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `[]`, body)
			},
		},
		{
			name:   "準正常系_blogIdが空の場合400を返す",
			blogId: "",
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("FetchCommentsByBlogId", "").Return(nil, services_comments.ErrInvalidBlogID)
			},
			wantErr: services_comments.ErrInvalidBlogID,
		},
		{
			// 取得の失敗を 404（コメント不在）として返さない
			name:   "異常系_Serviceがエラーを返す場合500を返す",
			blogId: "1",
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("FetchCommentsByBlogId", "1").Return(nil, errService)
			},
			wantErr: errService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/comments/blog/"+tt.blogId, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("blogId")
			c.SetParamValues(tt.blogId)

			mockCommentService := new(services_comments.MockCommentService)
			handler := NewCommentHandler(mockCommentService)

			tt.setupMock(mockCommentService)

			err := handler.FetchCommentsByBlogId(c)
			if tt.wantErr != nil {
				assertHandlerError(t, err, rec, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				tt.checkBody(t, rec.Body.String())
			}

			mockCommentService.AssertExpectations(t)
		})
	}
}
