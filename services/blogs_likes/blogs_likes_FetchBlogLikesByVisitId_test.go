package services_blogs_likes

import (
	"backend/models"
	repositories_blogs_likes "backend/repositories/blogs_likes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_FetchBlogLikesByVisitId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// モックデータの設定
	mockData := []models.BlogLikeData{
		{
			BlogId:  "1",
			VisitId: "1",
		},
		{
			BlogId:  "2",
			VisitId: "1",
		},
	}

	// モックの設定
	mockBlogLikeRepository.On("FetchBlogLikesByVisitId", "1").Return(mockData, nil)

	// 実行
	blogLikesData, err := blogLikeService.FetchBlogLikesByVisitId("1")

	// エラーチェック
	assert.NoError(t, err)
	assert.Len(t, blogLikesData, 2)
	assert.Equal(t, "1", blogLikesData[0].BlogId)

	// モックが期待通りに呼び出されたかを確認
	mockBlogLikeRepository.AssertExpectations(t)
}

func TestService_FetchBlogLikesByVisitId_InvalidVisitId(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// 実行
	blogLikesData, err := blogLikeService.FetchBlogLikesByVisitId("")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blogLikesData)
	assert.Contains(t, "visitId is empty", err.Error())

	// モックの呼び出しを確認
	mockBlogLikeRepository.AssertNotCalled(t, "FetchBlogLikesByVisitId", "")
}

func TestService_FetchBlogLikesByVisitId_NoData(t *testing.T) {
	// モックリポジトリをインスタンス化
	mockBlogLikeRepository := new(repositories_blogs_likes.MockBlogLikeRepository)
	blogLikeService := NewBlogLikeService(mockBlogLikeRepository)

	// モックの設定
	mockBlogLikeRepository.On("FetchBlogLikesByVisitId", "1").Return(nil, errors.New("no data"))

	// 実行
	blogLikesData, err := blogLikeService.FetchBlogLikesByVisitId("1")

	// エラーチェック
	assert.Error(t, err)
	assert.Nil(t, blogLikesData)
	assert.Contains(t, "no data", err.Error())

	// モックが期待通りに呼び出されたかを確認
	mockBlogLikeRepository.AssertExpectations(t)
}
