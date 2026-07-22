package repositories_blog_likes

import "backend/models"

// BlogLikeRepository はいいねのデータアクセスを担うリポジトリインターフェース。
type BlogLikeRepository interface {
	// FetchBlogLikesByVisitId は VisitId によっていいねデータを取得する。
	//
	// 引数:
	//   - visitId: 取得対象の訪問者ID
	//
	// 戻り値:
	//   - []models.BlogLikesData: 取得したいいねデータの一覧
	//   - error: 取得に失敗した場合のエラー
	FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error)
	// IsBlogLiked はいいねが存在するか確認する。
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
	//   - visitId: いいねした訪問者ID
	//
	// 戻り値:
	//   - *models.BlogLikesData: 作成したいいねデータ
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

// BlogLikeRepositoryImpl は BlogLikeRepository の実装。
type BlogLikeRepositoryImpl struct{}

// NewBlogLikeRepository は BlogLikeRepository を実装した BlogLikeRepositoryImpl のポインタを生成する。
//
// 戻り値:
//   - BlogLikeRepository: 生成したリポジトリ実装
func NewBlogLikeRepository() BlogLikeRepository {
	return &BlogLikeRepositoryImpl{}
}
