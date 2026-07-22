package repositories_blog_likes

import "backend/models"

// BlogLikeRepository はいいねのデータアクセスを担うリポジトリインターフェース。
type BlogLikeRepository interface {
	// FetchBlogLikesByVisitId は VisitId によっていいねデータを取得する。
	FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error)
	// IsBlogLiked はいいねが存在するか確認する。
	IsBlogLiked(blogId, visitId string) (bool, error)
	// CreateBlogLike はいいねデータを作成する。
	CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error)
	// DeleteBlogLike はいいねデータを削除する。
	DeleteBlogLike(blogId, visitId string) error
}

// BlogLikeRepositoryImpl は BlogLikeRepository の実装。
type BlogLikeRepositoryImpl struct{}

// NewBlogLikeRepository は BlogLikeRepository を実装した BlogLikeRepositoryImpl のポインタを生成する。
func NewBlogLikeRepository() BlogLikeRepository {
	return &BlogLikeRepositoryImpl{}
}
