package services_blogs_test

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	services_blogs "backend/services/blogs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogs(t *testing.T) {
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
	mockBlogRepository.On("FetchBlogs").Return(mockBlogData, nil)

	blogs, err := blogService.FetchBlogs()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.Len(t, blogs, 2)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchUsers_EmptyList(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := services_blogs.NewBlogService(mockBlogRepository)

	// ブログが存在しない場合
	mockBlogRepository.On("FetchBlogs").Return([]models.BlogData{}, nil)

	blogs, err := blogService.FetchBlogs()

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.Len(t, blogs, 0)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}
