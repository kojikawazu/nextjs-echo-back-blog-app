package services_comments

import (
	"backend/models"
	repositories_comments "backend/repositories/comments"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateComment(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	commentService := NewCommentService(mockCommentRepo)

	// 入力データ
	blogId := "1"
	guestUser := "guestUser1"
	comment := "comment1"

	// 期待されるコメントデータ
	expectedComment := models.CommentData{
		ID:        "1",
		BlogId:    "1",
		GuestUser: "guestUser1",
		Comment:   "comment1",
		CreatedAt: time.Now(),
	}

	// モックの設定
	mockCommentRepo.On("CreateComment", blogId, guestUser, comment).Return(&expectedComment, nil)

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, expectedComment.Comment, blog.Comment)

	// モックの期待通りの呼び出しを検証
	mockCommentRepo.AssertExpectations(t)
}

func TestService_CreateComment_InvalidBlogId(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	commentService := NewCommentService(mockCommentRepo)

	// 入力データ
	blogId := ""
	guestUser := "guestUser1"
	comment := "comment1"

	// モックの設定
	mockCommentRepo.On("CreateComment", blogId, guestUser, comment).Return(nil, errors.New("invalid blogId"))

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid blogId", err.Error())

	// モックの期待通りの呼び出しを検証
	mockCommentRepo.AssertNotCalled(t, "CreateComment", blogId, guestUser, comment)
}

func TestService_CreateComment_InvalidGuestUser(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	commentService := NewCommentService(mockCommentRepo)

	// 入力データ
	blogId := "1"
	guestUser := ""
	comment := "comment1"

	// モックの設定
	mockCommentRepo.On("CreateComment", blogId, guestUser, comment).Return(nil, errors.New("invalid guestUser"))

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid guestUser", err.Error())

	// モックの期待通りの呼び出しを検証
	mockCommentRepo.AssertNotCalled(t, "CreateComment", blogId, guestUser, comment)
}

func TestService_CreateComment_InvalidComment(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	commentService := NewCommentService(mockCommentRepo)

	// 入力データ
	blogId := "1"
	guestUser := "guestUser1"
	comment := ""

	// モックの設定
	mockCommentRepo.On("CreateComment", blogId, guestUser, comment).Return(nil, errors.New("invalid comment"))

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid comment", err.Error())

	// モックの期待通りの呼び出しを検証
	mockCommentRepo.AssertNotCalled(t, "CreateComment", blogId, guestUser, comment)
}

func TestService_CreateComment_NotCreate(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	commentService := NewCommentService(mockCommentRepo)

	// 入力データ
	blogId := "1"
	guestUser := "guestUser1"
	comment := "comment1"

	// モックの設定
	mockCommentRepo.On("CreateComment", blogId, guestUser, comment).Return(nil, errors.New("failed to create comment"))

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "failed to create comment", err.Error())

	// モックの期待通りの呼び出しを検証
	mockCommentRepo.AssertExpectations(t)
}
