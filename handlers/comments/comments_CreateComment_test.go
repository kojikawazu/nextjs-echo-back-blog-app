package handlers_comments

import (
	"backend/models"
	services_comments "backend/services/comments"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateComment(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"blogId":    "1",
		"guestUser": "guestUser1",
		"comment":   "comment1",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/comments/create", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの振る舞いを設定
	mockCommentService.On("CreateComment", "1", "guestUser1", "comment1").Return(&models.CommentData{
		ID:        "1",
		BlogId:    "1",
		GuestUser: "guestUser1",
		Comment:   "comment1",
	}, nil)

	// テストを実行
	err = handler.CreateComment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "guestUser1")

	// モックの呼び出しを確認
	mockCommentService.AssertExpectations(t)
}

func TestHandler_CreateComment_InvalidBlogId(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"blogId":    "",
		"guestUser": "guestUser1",
		"comment":   "comment1",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/comments/create", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの振る舞いを設定
	mockCommentService.On("CreateComment", "", "guestUser1", "comment1").Return(nil, errors.New("invalid blogId"))

	// テストを実行
	err = handler.CreateComment(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid blogId")

	// モックの呼び出しを確認
	mockCommentService.AssertExpectations(t)
}

func TestHandler_CreateComment_InvalidGuestUser(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"blogId":    "1",
		"guestUser": "",
		"comment":   "comment1",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/comments/create", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの振る舞いを設定
	mockCommentService.On("CreateComment", "1", "", "comment1").Return(nil, errors.New("invalid guestUser"))

	// テストを実行
	err = handler.CreateComment(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid guestUser")

	// モックの呼び出しを確認
	mockCommentService.AssertExpectations(t)
}

func TestHandler_CreateComment_InvalidComment(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"blogId":    "1",
		"guestUser": "guestUser1",
		"comment":   "",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/comments/create", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの振る舞いを設定
	mockCommentService.On("CreateComment", "1", "guestUser1", "").Return(nil, errors.New("invalid comment"))

	// テストを実行
	err = handler.CreateComment(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid comment")

	// モックの呼び出しを確認
	mockCommentService.AssertExpectations(t)
}

func TestHandler_CreateComment_NotCreate(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"blogId":    "1",
		"guestUser": "guestUser1",
		"comment":   "comment1",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/comments/create", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの振る舞いを設定
	mockCommentService.On("CreateComment", "1", "guestUser1", "comment1").Return(nil, errors.New("failed to create comment"))

	// テストを実行
	err = handler.CreateComment(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to create comment")

	// モックの呼び出しを確認
	mockCommentService.AssertExpectations(t)
}

func TestHandler_CreateComment_ServerError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()

	// JSONデータを作成
	requestBody := map[string]string{
		"blogId":    "1",
		"guestUser": "guestUser1",
		"comment":   "comment1",
	}

	// JSONデータをエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// リクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/comments/create", bytes.NewReader(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックサービスの生成
	mockCommentService := new(services_comments.MockCommentService)
	handler := NewCommentHandler(mockCommentService)

	// モックの振る舞いを設定
	mockCommentService.On("CreateComment", "1", "guestUser1", "comment1").Return(nil, errors.New("server error"))

	// テストを実行
	err = handler.CreateComment(c)
	assert.NoError(t, err)

	// ステータスコードとレスポンス内容の確認
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error creating comment")

	// モックの呼び出しを確認
	mockCommentService.AssertExpectations(t)
}
