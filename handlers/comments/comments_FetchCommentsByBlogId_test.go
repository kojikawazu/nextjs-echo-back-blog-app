package handlers_comments

import (
	"backend/models"
	services_comments "backend/services/comments"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_FetchCommentsByBlogId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	// パスパラメータとして blogId を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/comments/blog/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("blogId")
	c.SetParamValues("1")

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックデータの設定
	mockComment := []models.CommentData{
		{
			ID:        "1",
			BlogId:    "1",
			GuestUser: "guestUser1",
			Comment:   "comment1",
			CreatedAt: time.Now(),
		},
		{
			ID:        "2",
			BlogId:    "1",
			GuestUser: "guestUser2",
			Comment:   "comment2",
			CreatedAt: time.Now(),
		},
	}
	mockCommentService.On("FetchCommentsByBlogId", "1").Return(mockComment, nil)

	// ハンドラーを実行
	err := handler.FetchCommentsByBlogId(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "guestUser1")

	// モックが期待通りに呼び出されたかを確認
	mockCommentService.AssertExpectations(t)
}

func TestHandler_FetchCommentsByBlogId_InvalidBlogId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	// パスパラメータとして blogId を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/comments/blog/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("blogId")
	c.SetParamValues("1")

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの設定
	mockCommentService.On("FetchCommentsByBlogId", "1").Return(nil, errors.New("invalid blogId"))

	// ハンドラーを実行
	err := handler.FetchCommentsByBlogId(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid blogId")

	// モックが期待通りに呼び出されたかを確認
	mockCommentService.AssertExpectations(t)
}

func TestHandler_FetchCommentsByBlogId_NotComment(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	// パスパラメータとして blogId を指定する
	req := httptest.NewRequest(http.MethodGet, "/api/comments/blog/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// パスパラメータを設定
	c.SetParamNames("blogId")
	c.SetParamValues("1")

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの設定
	mockCommentService.On("FetchCommentsByBlogId", "1").Return(nil, errors.New("comments not found"))

	// ハンドラーを実行
	err := handler.FetchCommentsByBlogId(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Comments not found")

	// モックが期待通りに呼び出されたかを確認
	mockCommentService.AssertExpectations(t)
}
