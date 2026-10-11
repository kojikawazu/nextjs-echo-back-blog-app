package handlers_blog_comments

import (
	"backend/apperrors"
	"backend/models"
	services_comments "backend/services/blog_comments"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateComment(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(mockService *services_comments.MockCommentService)
		expectedStatus int
		expectedBody   string
		wantErr        error // 非 nil のときはハンドラがこのエラーを返すことを検証する
	}{
		{
			name:        "正常系_コメント作成成功",
			requestBody: map[string]string{"blogId": "1", "guestUser": "guestUser1", "comment": "comment1"},
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("CreateComment", "1", "guestUser1", "comment1").Return(&models.BlogCommentsData{
					ID:        "1",
					BlogId:    "1",
					GuestUser: "guestUser1",
					Comment:   "comment1",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"1","blog_id":"1","guest_user":"guestUser1","comment":"comment1","created_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:        "準正常系_blogIdが空文字の場合400を返す",
			requestBody: map[string]string{"blogId": "", "guestUser": "guestUser1", "comment": "comment1"},
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("CreateComment", "", "guestUser1", "comment1").Return(nil, services_comments.ErrInvalidBlogID)
			},
			wantErr: services_comments.ErrInvalidBlogID,
		},
		{
			name:        "準正常系_guestUserが空文字の場合400を返す",
			requestBody: map[string]string{"blogId": "1", "guestUser": "", "comment": "comment1"},
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("CreateComment", "1", "", "comment1").Return(nil, services_comments.ErrInvalidGuestUser)
			},
			wantErr: services_comments.ErrInvalidGuestUser,
		},
		{
			name:        "準正常系_commentが空文字の場合400を返す",
			requestBody: map[string]string{"blogId": "1", "guestUser": "guestUser1", "comment": ""},
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("CreateComment", "1", "guestUser1", "").Return(nil, services_comments.ErrInvalidComment)
			},
			wantErr: services_comments.ErrInvalidComment,
		},
		{
			name:        "異常系_Serviceがエラーを返す場合500を返す",
			requestBody: map[string]string{"blogId": "1", "guestUser": "guestUser1", "comment": "comment1"},
			setupMock: func(mockService *services_comments.MockCommentService) {
				mockService.On("CreateComment", "1", "guestUser1", "comment1").Return(nil, errService)
			},
			wantErr: errService,
		},
		{
			name:        "異常系_不正なJSONボディの場合400を返す",
			requestBody: nil,
			setupMock:   func(mockService *services_comments.MockCommentService) {},
			wantErr:     apperrors.ErrInvalidRequestBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			var body *bytes.Buffer
			if tt.requestBody == nil {
				body = bytes.NewBufferString("{invalid json}")
			} else {
				jsonData, err := json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				body = bytes.NewBuffer(jsonData)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/comments/create", body)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockCommentService := new(services_comments.MockCommentService)
			handler := NewCommentHandler(mockCommentService)

			tt.setupMock(mockCommentService)

			err := handler.CreateComment(c)
			if tt.wantErr != nil {
				assertHandlerError(t, err, rec, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			}

			mockCommentService.AssertExpectations(t)
		})
	}
}
