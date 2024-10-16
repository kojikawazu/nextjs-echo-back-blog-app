package services_comments

import (
	"backend/models"
	repositories_comments "backend/repositories/comments"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchCommentsByBlogId(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	serviceComment := NewCommentService(mockCommentRepo)

	// テストデータ
	blogId := "1"

	mockCommentData := []models.CommentData{
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

	// モックリポジトリのメソッドをモック化
	mockCommentRepo.On("FetchCommentsByBlogId", blogId).Return(mockCommentData, nil)

	// テスト対象メソッドの実行
	comments, err := serviceComment.FetchCommentsByBlogId(blogId)

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.Len(t, comments, 2)

	// モックが期待通りに呼び出されたかを確認
	mockCommentRepo.AssertExpectations(t)
}

func TestService_FetchCommentsByBlogId_InvalidBlogId(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	commentService := NewCommentService(mockCommentRepo)

	// テストデータ
	blogId := ""

	// テスト対象メソッドの実行
	comments, err := commentService.FetchCommentsByBlogId(blogId)

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, comments)
	assert.Equal(t, "invalid blogId", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockCommentRepo.AssertNotCalled(t, "FetchCommentsByBlogId", blogId)
}

func TestService_FetchCommentsByBlogId_NotComments(t *testing.T) {
	// モックリポジトリの生成
	mockCommentRepo := new(repositories_comments.MockCommentRepository)
	commentService := NewCommentService(mockCommentRepo)

	// テストデータ
	blogId := "1"

	// モックリポジトリのメソッドをモック化
	mockCommentRepo.On("FetchCommentsByBlogId", blogId).Return(nil, errors.New("comments not found"))

	// テスト対象メソッドの実行
	comments, err := commentService.FetchCommentsByBlogId(blogId)

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, comments)
	assert.Equal(t, "comments not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockCommentRepo.AssertExpectations(t)
}
