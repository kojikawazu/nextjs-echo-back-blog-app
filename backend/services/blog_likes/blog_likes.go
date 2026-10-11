package services_blog_likes

import (
	"backend/models"
	"fmt"
	"log"
)

// FetchBlogLikesByVisitId は VisitId に紐づくいいねデータを取得する。
//
// 引数:
//   - visitId: 取得対象の訪問者ID
//
// 戻り値:
//   - []models.BlogLikesData: 取得したいいねデータの一覧
//   - error: visitId が空の場合は ErrVisitIDEmpty、その他取得に失敗した場合のエラー
func (s *BlogLikeServiceImpl) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error) {
	log.Println("FetchBlogLikesByVisitId start...")

	// バリデーション
	if visitId == "" {
		log.Println("VisitId is empty")
		return nil, ErrVisitIDEmpty
	}

	log.Println("validation passed")

	// いいねデータを取得
	blogLikes, err := s.BlogLikeRepository.FetchBlogLikesByVisitId(visitId)
	if err != nil {
		return nil, fmt.Errorf("fetch blog likes: %w", err)
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
//   - error: blogId・visitId が空の場合は ErrBlogIDOrVisitIDEmpty、その他確認に失敗した場合のエラー
func (s *BlogLikeServiceImpl) IsBlogLiked(blogId, visitId string) (bool, error) {
	log.Println("IsBlogLiked start...")

	// バリデーション
	if blogId == "" || visitId == "" {
		log.Println("BlogId or VisitId is empty")
		return false, ErrBlogIDOrVisitIDEmpty
	}
	log.Println("validation passed")

	// いいねデータが存在するか確認
	isLiked, err := s.BlogLikeRepository.IsBlogLiked(blogId, visitId)
	if err != nil {
		return false, fmt.Errorf("check blog liked: %w", err)
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
//   - error: blogId・visitId が空の場合は ErrBlogIDOrVisitIDEmpty、いいね済みの場合は ErrAlreadyLiked、その他作成に失敗した場合のエラー
func (s *BlogLikeServiceImpl) CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error) {
	log.Println("CreateBlogLike start...")

	// バリデーション
	if blogId == "" || visitId == "" {
		log.Println("BlogId or VisitId is empty")
		return nil, ErrBlogIDOrVisitIDEmpty
	}
	// ブログがいいねされているか確認
	isLiked, err := s.BlogLikeRepository.IsBlogLiked(blogId, visitId)
	if err != nil {
		return nil, fmt.Errorf("check blog liked: %w", err)
	}
	if isLiked {
		log.Println("Blog is already liked")
		return nil, ErrAlreadyLiked
	}

	log.Println("validation passed")

	// いいねデータを作成
	blogLike, err := s.BlogLikeRepository.CreateBlogLike(blogId, visitId)
	if err != nil {
		return nil, fmt.Errorf("create blog like: %w", err)
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
		return fmt.Errorf("delete blog like: %w", err)
	}

	return nil
}
