package services_blogs_likes

import (
	"backend/models"
	repositories_blogs_likes "backend/repositories/blogs_likes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_CreateBlogLike(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// モックデータ
	blogLikeData := &models.BlogLikeData{
		BlogId:  "1",
		VisitId: "1",
	}

	// モックの設定
	mockBlogLikeRepository.On("IsBlogLiked", "1", "1").Return(false, nil)
	mockBlogLikeRepository.On("CreateBlogLike", "1", "1").Return(blogLikeData, nil)

	// 実行
	createdBlogLikeData, err := blogLikeService.CreateBlogLike("1", "1")

	// エラーチェック
	assert.NoError(t, err)
	assert.Equal(t, blogLikeData, createdBlogLikeData)

	// モックが期待通りに呼び出されたかを確認
	mockBlogLikeRepository.AssertExpectations(t)
}

func TestService_CreateBlogLike_InValidBlogId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// 実行
	createdBlogLikeData, err := blogLikeService.CreateBlogLike("", "1")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, createdBlogLikeData)
	assert.Contains(t, "blogId or VisitId is empty", err.Error())

	// モックの呼び出し確認
	mockBlogLikeRepository.AssertNotCalled(t, "IsBlogLiked", "", "1")
	mockBlogLikeRepository.AssertNotCalled(t, "CreateBlogLike", "", "1")
}

func TestService_CreateBlogLike_InValidVisitId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// 実行
	createdBlogLikeData, err := blogLikeService.CreateBlogLike("1", "")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, createdBlogLikeData)
	assert.Contains(t, "blogId or VisitId is empty", err.Error())

	// モックの呼び出し確認
	mockBlogLikeRepository.AssertNotCalled(t, "IsBlogLiked", "1", "")
	mockBlogLikeRepository.AssertNotCalled(t, "CreateBlogLike", "1", "")
}

func TestService_CreateBlogLike_AlreadyBlogLike(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// モックの設定
	mockBlogLikeRepository.On("IsBlogLiked", "1", "1").Return(true, nil)

	// 実行
	createdBlogLikeData, err := blogLikeService.CreateBlogLike("1", "1")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, createdBlogLikeData)
	assert.Contains(t, "blog is already liked", err.Error())

	// モックの呼び出し確認
	mockBlogLikeRepository.AssertExpectations(t)
	mockBlogLikeRepository.AssertNotCalled(t, "CreateBlogLike", "1", "1")
}

func TestService_CreateBlogLike_NotCreate(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// モックの設定
	mockBlogLikeRepository.On("IsBlogLiked", "1", "1").Return(false, nil)
	mockBlogLikeRepository.On("CreateBlogLike", "1", "1").Return(nil, errors.New("not created"))

	// 実行
	createdBlogLikeData, err := blogLikeService.CreateBlogLike("1", "1")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, createdBlogLikeData)
	assert.Contains(t, "not created", err.Error())

	// モックの呼び出し確認
	mockBlogLikeRepository.AssertExpectations(t)
}
