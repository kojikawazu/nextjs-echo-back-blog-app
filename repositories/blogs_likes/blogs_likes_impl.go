package repositories_blogs_likes

import "backend/models"

// BlogLikeRepositoryインターフェース
type BlogLikeRepository interface {
	IsBlogLiked(blogId, visitId string) (bool, error)
	CreateBlogLike(blogId, visitId string) (*models.BlogLikeData, error)
	DeleteBlogLike(blogId, visitId string) error
}

type BlogLikeRepositoryImpl struct{}

// BlogLikeRepositoryインターフェースを実装したBlogLikeRepositoryImplのポインタを返す
func NewBlogLikeRepository() BlogLikeRepository {
	return &BlogLikeRepositoryImpl{}
}
