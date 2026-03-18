package services_blog_likes

import (
	"backend/models"
	repositories_blogs_likes "backend/repositories/blog_likes"
)

// BlogLikeServiceインターフェース
type BlogLikeService interface {
	FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error)
	IsBlogLiked(blogId, visitId string) (bool, error)
	CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error)
	DeleteBlogLike(blogId, visitId string) error
}

type BlogLikeServiceImpl struct {
	BlogLikeRepository repositories_blogs_likes.BlogLikeRepository
}

// BlogLikeServiceインターフェースを実装したBlogLikeServiceImplのポインタを返す
func NewBlogLikeService(
	blogLikeRepository repositories_blogs_likes.BlogLikeRepository,
) BlogLikeService {
	return &BlogLikeServiceImpl{
		BlogLikeRepository: blogLikeRepository,
	}
}
