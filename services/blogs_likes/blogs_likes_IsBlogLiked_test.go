package services_blogs_likes

import (
	repositories_blogs_likes "backend/repositories/blogs_likes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_IsBlogLiked(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// モックの設定
	mockBlogLikeRepository.On("IsBlogLiked", "1", "1").Return(true, nil)

	// 実行
	isLiked, err := blogLikeService.IsBlogLiked("1", "1")

	// エラーチェック
	assert.NoError(t, err)
	assert.True(t, isLiked)

	// モックが期待通りに呼び出されたかを確認
	mockBlogLikeRepository.AssertExpectations(t)
}

func TestService_IsBlogLiked_InvalidBlogId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// 実行
	isLiked, err := blogLikeService.IsBlogLiked("", "1")

	// エラーチェック
	assert.Error(t, err)
	assert.False(t, isLiked)
	assert.Contains(t, "BlogId or VisitId is empty", err.Error())

	// モックが期待通りに呼び出されてないかを確認
	mockBlogLikeRepository.AssertNotCalled(t, "IsBlogLiked", "", "1")
}

func TestService_IsBlogLiked_InvalidVisitId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// 実行
	isLiked, err := blogLikeService.IsBlogLiked("1", "")

	// エラーチェック
	assert.Error(t, err)
	assert.False(t, isLiked)
	assert.Contains(t, "BlogId or VisitId is empty", err.Error())

	// モックが期待通りに呼び出されてないかを確認
	mockBlogLikeRepository.AssertNotCalled(t, "IsBlogLiked", "1", "")
}

func TestService_IsBlogLiked_NotLike(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// モックの設定
	mockBlogLikeRepository.On("IsBlogLiked", "1", "1").Return(false, errors.New("not found"))

	// 実行
	isLiked, err := blogLikeService.IsBlogLiked("1", "1")

	// エラーチェック
	assert.Error(t, err)
	assert.False(t, isLiked)
	assert.Contains(t, "not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogLikeRepository.AssertExpectations(t)
}
