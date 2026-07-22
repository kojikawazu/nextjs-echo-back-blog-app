package services_blog_likes

import (
	"backend/models"
	repositories_blogs_likes "backend/repositories/blog_likes"
)

// BlogLikeService はいいねに関する処理を提供するサービスインターフェース。
type BlogLikeService interface {
	// FetchBlogLikesByVisitId は VisitId に紐づくいいねデータを取得する。
	//
	// 引数:
	//   - visitId: 取得対象の訪問者ID
	//
	// 戻り値:
	//   - []models.BlogLikesData: 取得したいいねデータの一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error)
	// IsBlogLiked はいいねが存在するかを確認する。
	//
	// 引数:
	//   - blogId: 確認対象のブログID
	//   - visitId: 確認対象の訪問者ID
	//
	// 戻り値:
	//   - bool: いいねが存在する場合は true
	//   - error: 確認に失敗した場合のエラー
	IsBlogLiked(blogId, visitId string) (bool, error)
	// CreateBlogLike はいいねデータを作成する。
	//
	// 引数:
	//   - blogId: いいね対象のブログID
	//   - visitId: いいねする訪問者ID
	//
	// 戻り値:
	//   - *models.BlogLikesData: 作成されたいいねデータ
	//   - error: 作成に失敗した場合のエラー
	CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error)
	// DeleteBlogLike はいいねデータを削除する。
	//
	// 引数:
	//   - blogId: 削除対象のブログID
	//   - visitId: 削除対象の訪問者ID
	//
	// 戻り値:
	//   - error: 削除に失敗した場合のエラー
	DeleteBlogLike(blogId, visitId string) error
}

// BlogLikeServiceImpl は BlogLikeService の実装。
type BlogLikeServiceImpl struct {
	BlogLikeRepository repositories_blogs_likes.BlogLikeRepository
}

// NewBlogLikeService は BlogLikeService インターフェースを実装した BlogLikeServiceImpl を生成する。
//
// 引数:
//   - blogLikeRepository: いいねデータへのアクセスを担うリポジトリ
//
// 戻り値:
//   - BlogLikeService: 生成されたいいねサービス
func NewBlogLikeService(
	blogLikeRepository repositories_blogs_likes.BlogLikeRepository,
) BlogLikeService {
	return &BlogLikeServiceImpl{
		BlogLikeRepository: blogLikeRepository,
	}
}
