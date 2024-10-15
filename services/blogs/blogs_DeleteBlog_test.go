package services_blogs

import (
	repositories_blogs "backend/repositories/blogs"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_DeleteBlog(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"

	// モックの設定
	mockBlogRepository.On("DeleteBlog", id).Return(nil, nil)

	// テスト対象メソッドの呼び出し
	err := blogService.DeleteBlog(id)

	// アサーション
	assert.NoError(t, err)
	assert.Nil(t, err)

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertExpectations(t)
}

func TestService_DeleteBlog_InvalidId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := ""

	// モックの設定
	mockBlogRepository.On("DeleteBlog", id).Return(nil, errors.New("invalid id"))

	// テスト対象メソッドの呼び出し
	err := blogService.DeleteBlog(id)

	// アサーション
	assert.Error(t, err)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid id", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertNotCalled(t, "DeleteBlog", id)
}

func TestService_DeleteBlog_NotBlog(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"

	// モックの設定
	mockBlogRepository.On("DeleteBlog", id).Return(errors.New("failed to delete blog"))

	// テスト対象メソッドの呼び出し
	err := blogService.DeleteBlog(id)

	// アサーション
	assert.Error(t, err)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to delete blog", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertExpectations(t)
}
