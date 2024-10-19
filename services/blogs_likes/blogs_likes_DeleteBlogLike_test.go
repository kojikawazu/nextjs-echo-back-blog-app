package services_blogs_likes

import (
	repositories_blogs_likes "backend/repositories/blogs_likes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_DeleteBlogLike(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	mockBlogLikeRepository.On("DeleteBlogLike", "1", "1").Return(nil)

	// 実行
	err := blogLikeService.DeleteBlogLike("1", "1")

	// エラーチェック
	assert.NoError(t, err)

	// モックが期待通りに呼び出されたかを確認
	mockBlogLikeRepository.AssertExpectations(t)
}

func TestService_DeleteBlogLike_NoDelete(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	mockBlogLikeRepository.On("DeleteBlogLike", "1", "1").Return(errors.New("Delete Error"))

	// 実行
	err := blogLikeService.DeleteBlogLike("1", "1")

	// エラーチェック
	assert.Error(t, err)
	assert.Contains(t, "Delete Error", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogLikeRepository.AssertExpectations(t)
}
