package services_blogs_test

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_FetchBlogPopular(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	mockBlogData := []models.BlogData{
		{
			ID:        "1",
			UserId:    "1",
			Title:     "title1",
			GithubUrl: "https://github.com/user/repo1",
			Category:  "Category1",
			Tags:      "Tag1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			UserId:    "2",
			Title:     "title2",
			GithubUrl: "https://github.com/user/repo2",
			Category:  "Category2",
			Tags:      "Tag2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogPopular", 2).Return(mockBlogData, nil)

	blogData, err := blogService.FetchBlogPopular(2)

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blogData)
	assert.Len(t, blogData, 2)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogPopular_InvalidCount(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogPopular", 0).Return(nil, errors.New("No data"))

	blogData, err := blogService.FetchBlogPopular(0)

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blogData)

	// モックの呼び出しがないことを確認
	mockBlogRepository.AssertNotCalled(t, "FetchBlogPopular", mock.Anything)
}
