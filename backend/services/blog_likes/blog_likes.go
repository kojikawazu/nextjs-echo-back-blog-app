package services_blog_likes

import (
	"backend/models"
	"errors"
	"log"
)

// FetchBlogLikesByVisitId は VisitId に紐づくいいねデータを取得する。
//
// 引数:
//   - visitId: 取得対象の訪問者ID
//
// 戻り値:
//   - []models.BlogLikesData: 取得したいいねデータの一覧
//   - error: 取得に失敗した場合のエラー
func (s *BlogLikeServiceImpl) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error) {
	log.Println("FetchBlogLikesByVisitId start...")

	// バリデーション
	if visitId == "" {
		log.Println("VisitId is empty")
		return nil, errors.New("visitId is empty")
	}

	log.Println("validation passed")

	// いいねデータを取得
	blogLikes, err := s.BlogLikeRepository.FetchBlogLikesByVisitId(visitId)
	if err != nil {
		return nil, err
	}

	return blogLikes, nil
}

// IsBlogLiked はいいねが存在するかを確認する。
//
// 引数:
//   - blogId: 確認対象のブログID
//   - visitId: 確認対象の訪問者ID
//
// 戻り値:
//   - bool: いいねが存在する場合は true
//   - error: 確認に失敗した場合のエラー
func (s *BlogLikeServiceImpl) IsBlogLiked(blogId, visitId string) (bool, error) {
	log.Println("IsBlogLiked start...")

	// バリデーション
	if blogId == "" || visitId == "" {
		log.Println("BlogId or VisitId is empty")
		return false, errors.New("BlogId or VisitId is empty")
	}
	log.Println("validation passed")

	// いいねデータが存在するか確認
	isLiked, err := s.BlogLikeRepository.IsBlogLiked(blogId, visitId)
	if err != nil {
		return false, err
	}

	return isLiked, nil
}

// CreateBlogLike はいいねデータを作成する。
//
// 引数:
//   - blogId: いいね対象のブログID
//   - visitId: いいねする訪問者ID
//
// 戻り値:
//   - *models.BlogLikesData: 作成されたいいねデータ
//   - error: 作成に失敗した場合のエラー
func (s *BlogLikeServiceImpl) CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error) {
	log.Println("CreateBlogLike start...")

	// バリデーション
	if blogId == "" || visitId == "" {
		log.Println("BlogId or VisitId is empty")
		return nil, errors.New("blogId or VisitId is empty")
	}
	// ブログがいいねされているか確認
	IsBlogLiked, _ := s.BlogLikeRepository.IsBlogLiked(blogId, visitId)
	if IsBlogLiked {
		log.Println("Blog is already liked")
		return nil, errors.New("blog is already liked")
	}

	log.Println("validation passed")

	// いいねデータを作成
	blogLike, err := s.BlogLikeRepository.CreateBlogLike(blogId, visitId)
	if err != nil {
		return nil, err
	}

	return blogLike, nil
}

// DeleteBlogLike はいいねデータを削除する。
//
// 引数:
//   - blogId: 削除対象のブログID
//   - visitId: 削除対象の訪問者ID
//
// 戻り値:
//   - error: 削除に失敗した場合のエラー
func (s *BlogLikeServiceImpl) DeleteBlogLike(blogId, visitId string) error {
	log.Println("DeleteBlogLike start...")

	// いいねデータを削除
	err := s.BlogLikeRepository.DeleteBlogLike(blogId, visitId)
	if err != nil {
		return err
	}

	return nil
}
