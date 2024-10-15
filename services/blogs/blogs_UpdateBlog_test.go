package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestService_UpdateBlog(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"
	userId := "user1"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// 期待されるブログデータ
	expectedBlog := models.BlogData{
		ID:          "123",
		UserId:      userId,
		Title:       title,
		GithubUrl:   githubURL,
		Category:    category,
		Description: description,
		Tags:        tags,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// モックの設定
	mockBlogRepository.On("UpdateBlog", id, title, githubURL, category, description, tags).Return(&expectedBlog, nil)

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.NoError(t, err)
	assert.NotNil(t, blog)
	assert.Equal(t, &expectedBlog, blog)

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertExpectations(t)
}

func TestService_UpdateBlog_InvalidId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := ""
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid id", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertNotCalled(t, "UpdateBlog", id, title, githubURL, category, description, tags)
}

func TestService_UpdateBlog_InvalidTitle(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"
	title := ""
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid title", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertNotCalled(t, "UpdateBlog", id, title, githubURL, category, description, tags)
}

func TestService_UpdateBlog_InvalidGithubUrl(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"
	title := "Test Blog"
	githubURL := ""
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid githubUrl", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertNotCalled(t, "UpdateBlog", id, title, githubURL, category, description, tags)
}

func TestService_UpdateBlog_InvalidCategory(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := ""
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid category", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertNotCalled(t, "UpdateBlog", id, title, githubURL, category, description, tags)
}

func TestService_UpdateBlog_InvalidDescription(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := ""
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid description", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertNotCalled(t, "UpdateBlog", id, title, githubURL, category, description, tags)
}

func TestService_UpdateBlog_InvalidTags(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := ""

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "invalid tags", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertNotCalled(t, "UpdateBlog", id, title, githubURL, category, description, tags)
}

func TestService_UpdateBlog_NoUpdate(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	id := "123"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// モックの設定
	mockBlogRepository.On("UpdateBlog", id, title, githubURL, category, description, tags).Return(nil, errors.New("no update"))

	// テスト対象メソッドの呼び出し
	blog, err := blogService.UpdateBlog(id, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "failed to update blog", err.Error())

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertExpectations(t)
}
