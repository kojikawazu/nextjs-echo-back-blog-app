package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogs(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	mockBlogData := []models.BlogData{
		{
			ID:        "1",
			UserId:    "1",
			Title:     "title1",
			GithubUrl: "https://github.com/user/repo1",
			Category:  "Category1",
			Tag:       "Tag1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			UserId:    "2",
			Title:     "title2",
			GithubUrl: "https://github.com/user/repo2",
			Category:  "Category2",
			Tag:       "Tag2",
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
	blogService := NewBlogService(mockBlogRepository)

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

func TestService_FetchBlogByUserId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// ブログが存在する場合
	mockBlogData := &models.BlogData{
		ID:        "1",
		UserId:    "1",
		Title:     "title1",
		GithubUrl: "https://github.com/user/repo1",
		Category:  "Category1",
		Tag:       "Tag1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockBlogRepository.On("FetchBlogByUserId", "1").Return(mockBlogData, nil)

	// サービス層メソッドの実行
	blog, err := blogService.FetchBlogByUserId("1")

	// エラーチェック
	assert.NoError(t, err)

	// データが期待通りか確認
	assert.NotNil(t, blog) // blogがnilでないことを確認
	assert.Equal(t, "title1", blog.Title)

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogByUserId_InvalidCases(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// サービス層メソッドの実行
	_, err := blogService.FetchBlogByUserId("")

	// エラーチェック
	assert.Error(t, err)
	assert.Equal(t, "invalid userId", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}

func TestService_FetchBlogByUserId_NotUser(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// "blog not found" エラーメッセージを返すように設定
	mockBlogRepository.On("FetchBlogByUserId", "2").Return(nil, errors.New("blog not found"))

	// サービス層メソッドの実行
	blog, err := blogService.FetchBlogByUserId("2")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "blog not found", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogRepository.AssertExpectations(t)
}
