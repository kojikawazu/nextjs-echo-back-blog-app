package services_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogCategories(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	mockBlogCategories := []string{
		"Category1",
		"Category2",
	}

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogCategories").Return(mockBlogCategories, nil)

	blogCategories, err := blogService.FetchBlogCategories()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blogCategories)
	assert.Len(t, blogCategories, 2)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogCategories_NoData(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	// モックデータ
	mockBlogCategories := []string{}

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogCategories").Return(mockBlogCategories, nil)

	blogCategories, err := blogService.FetchBlogCategories()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blogCategories)
	assert.Len(t, blogCategories, 0)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogCategories_ErrorCase(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogCategories").Return(nil, errors.New("No data"))

	blogCategories, err := blogService.FetchBlogCategories()

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blogCategories)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}
