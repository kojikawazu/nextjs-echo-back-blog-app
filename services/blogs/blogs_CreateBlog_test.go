package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateBlog(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
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
	mockBlogRepository.On("CreateBlog", userId, title, githubURL, category, description, tags).Return(expectedBlog, nil)

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.NoError(t, err)
	assert.Equal(t, expectedBlog, blog)

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertExpectations(t)
}

func TestService_CreateBlog_InvalidUserId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	userId := ""
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "invalid userId", err.Error())
	assert.Equal(t, models.BlogData{}, blog)

	// モックの呼び出しがないことを確認
	mockBlogRepository.AssertNotCalled(t, "CreateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateBlog_InvalidTitle(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	userId := "user1"
	title := ""
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "invalid title", err.Error())
	assert.Equal(t, models.BlogData{}, blog)

	// モックの呼び出しがないことを確認
	mockBlogRepository.AssertNotCalled(t, "CreateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateBlog_InvalidGitHubURL(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	userId := "user1"
	title := "Test Blog"
	githubURL := ""
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "invalid githubUrl", err.Error())
	assert.Equal(t, models.BlogData{}, blog)

	// モックの呼び出しがないことを確認
	mockBlogRepository.AssertNotCalled(t, "CreateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateBlog_InvalidCategory(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	userId := "user1"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := ""
	description := "This is a test blog."
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "invalid category", err.Error())
	assert.Equal(t, models.BlogData{}, blog)

	// モックの呼び出しがないことを確認
	mockBlogRepository.AssertNotCalled(t, "CreateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateBlog_InvalidDescription(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	userId := "user1"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := ""
	tags := "go, testing"

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "invalid description", err.Error())
	assert.Equal(t, models.BlogData{}, blog)

	// モックの呼び出しがないことを確認
	mockBlogRepository.AssertNotCalled(t, "CreateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateBlog_InvalidTags(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	userId := "user1"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := ""

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "invalid tags", err.Error())
	assert.Equal(t, models.BlogData{}, blog)

	// モックの呼び出しがないことを確認
	mockBlogRepository.AssertNotCalled(t, "CreateBlog", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestService_CreateBlog_RepositoryError(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogRepository := new(repositories_blogs.MockBlogRepository)
	blogService := NewBlogService(mockBlogRepository)

	// 入力データ
	userId := "user1"
	title := "Test Blog"
	githubURL := "https://github.com/user/repo"
	category := "Tech"
	description := "This is a test blog."
	tags := "go, testing"

	// モックの設定: リポジトリがエラーを返す
	mockBlogRepository.On("CreateBlog", userId, title, githubURL, category, description, tags).Return(models.BlogData{}, errors.New("repository failure"))

	// テスト対象メソッドの呼び出し
	blog, err := blogService.CreateBlog(userId, title, githubURL, category, description, tags)

	// アサーション
	assert.Error(t, err)
	assert.Equal(t, "failed to create blog", err.Error())
	assert.Equal(t, models.BlogData{}, blog)

	// モックの期待通りの呼び出しを検証
	mockBlogRepository.AssertExpectations(t)
}
