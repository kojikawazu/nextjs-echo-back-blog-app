package repositories_blog_likes

import "backend/models"

// BlogLikeRepositoryインターフェース
type BlogLikeRepository interface {
	FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error)
	IsBlogLiked(blogId, visitId string) (bool, error)
	CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error)
	DeleteBlogLike(blogId, visitId string) error
}

type BlogLikeRepositoryImpl struct{}

// BlogLikeRepositoryインターフェースを実装したBlogLikeRepositoryImplのポインタを返す
func NewBlogLikeRepository() BlogLikeRepository {
	return &BlogLikeRepositoryImpl{}
}
