package services_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogTags(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	mockBlogTags := []string{
		"Tag1",
		"Tag2",
	}

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogTags").Return(mockBlogTags, nil)

	blogTags, err := blogService.FetchBlogTags()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blogTags)
	assert.Len(t, blogTags, 2)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogTags_NoData(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	// モックデータ
	mockBlogTags := []string{}

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogTags").Return(mockBlogTags, nil)

	blogTags, err := blogService.FetchBlogTags()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blogTags)
	assert.Len(t, blogTags, 0)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogTags_ErrorCase(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogTags").Return(nil, errors.New("No data"))

	blogTags, err := blogService.FetchBlogTags()

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blogTags)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}
