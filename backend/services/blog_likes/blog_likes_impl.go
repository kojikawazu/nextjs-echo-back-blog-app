package services_blog_likes

import (
	"backend/models"
	repositories_blogs_likes "backend/repositories/blog_likes"
)

// BlogLikeService はいいねに関する処理を提供するサービスインターフェース。
type BlogLikeService interface {
	// FetchBlogLikesByVisitId は VisitId に紐づくいいねデータを取得する。
	FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error)
	// IsBlogLiked はいいねが存在するかを確認する。
	IsBlogLiked(blogId, visitId string) (bool, error)
	// CreateBlogLike はいいねデータを作成する。
	CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error)
	// DeleteBlogLike はいいねデータを削除する。
	DeleteBlogLike(blogId, visitId string) error
}

// BlogLikeServiceImpl は BlogLikeService の実装。
type BlogLikeServiceImpl struct {
	BlogLikeRepository repositories_blogs_likes.BlogLikeRepository
}

// NewBlogLikeService は BlogLikeService インターフェースを実装した BlogLikeServiceImpl を生成する。
func NewBlogLikeService(
	blogLikeRepository repositories_blogs_likes.BlogLikeRepository,
) BlogLikeService {
	return &BlogLikeServiceImpl{
		BlogLikeRepository: blogLikeRepository,
	}
}
