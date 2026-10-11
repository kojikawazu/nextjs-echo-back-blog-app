package services_blog_comments

import (
	"backend/models"
	repositories_comments "backend/repositories/blog_comments"
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
	expectedComment := models.BlogCommentsData{
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

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.ErrorIs(t, err, ErrInvalidBlogID)

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

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.ErrorIs(t, err, ErrInvalidGuestUser)

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

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.ErrorIs(t, err, ErrInvalidComment)

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
	repoErr := errors.New("db error")
	mockCommentRepo.On("CreateComment", blogId, guestUser, comment).Return(nil, repoErr)

	// テスト対象メソッドの呼び出し
	blog, err := commentService.CreateComment(blogId, guestUser, comment)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	// 原因（リポジトリのエラー）を辿れ、入力エラーの sentinel には該当しない（= 500 になる）
	assert.ErrorIs(t, err, repoErr)
	assert.NotErrorIs(t, err, ErrInvalidBlogID)

	// モックの期待通りの呼び出しを検証
	mockCommentRepo.AssertExpectations(t)
}
