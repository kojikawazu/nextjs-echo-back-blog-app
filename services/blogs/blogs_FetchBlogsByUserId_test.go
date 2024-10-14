package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogsByUserId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// ブログが存在する場合
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
	mockBlogRepository.On("FetchBlogsByUserId", "1").Return(mockBlogData, nil)

	// サービス層メソッドの実行
	blogs, err := blogService.FetchBlogsByUserId("1")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, blogs)
	assert.Len(t, blogs, 2)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogsByUserId_InvalidCases(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// サービス層メソッドの実行
	_, err := blogService.FetchBlogsByUserId("")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "invalid userId", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogsByUserId_NotUser(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// "blog not found" エラーメッセージを返すように設定
	mockBlogRepository.On("FetchBlogsByUserId", "2").Return(nil, errors.New("blog not found"))

	// サービス層メソッドの実行
	blogs, err := blogService.FetchBlogsByUserId("2")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blogs)
	assert.Equal(t, "blogs not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}
