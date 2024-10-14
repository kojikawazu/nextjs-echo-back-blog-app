package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogById(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// モックデータ
	mockBlogData := &models.BlogData{
		ID:        "1",
		UserId:    "1",
		Title:     "title1",
		GithubUrl: "https://github.com/user/repo1",
		Category:  "Category1",
		Tags:      "Tag1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ブログが存在する場合
	mockBlogRepository.On("FetchBlogById", "1").Return(mockBlogData, nil)

	// ブログデータを取得
	blog, err := blogService.FetchBlogById("1")

	// エラーチェック
	assert.NoError(t, err)
	assert.NotNil(t, blog)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogById_InvalidId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// IDが空の場合
	blog, err := blogService.FetchBlogById("")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid id", err.Error())

	// モックが呼ばれなかったことを確認
	mockBlogRepository.AssertNotCalled(t, "FetchBlogById")
}

func TestService_FetchBlogById_NotBlog(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// モックを設定
	mockBlogRepository.On("FetchBlogById", "1").Return(nil, errors.New("blog not found"))

	// ブログが存在しない場合
	blog, err := blogService.FetchBlogById("1")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "blog not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}
