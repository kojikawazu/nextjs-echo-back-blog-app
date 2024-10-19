package services_blogs_likes

import (
	"backend/models"
	"errors"
	"log"
)

// VisitIdに紐づくいいねデータを取得
func (s *BlogLikeServiceImpl) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikeData, error) {
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

// いいね存在するか確認
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

// いいねデータの作成
func (s *BlogLikeServiceImpl) CreateBlogLike(blogId, visitId string) (*models.BlogLikeData, error) {
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

// いいねデータの削除
func (s *BlogLikeServiceImpl) DeleteBlogLike(blogId, visitId string) error {
	log.Println("DeleteBlogLike start...")

	// いいねデータを削除
	err := s.BlogLikeRepository.DeleteBlogLike(blogId, visitId)
	if err != nil {
		return err
	}

	return nil
}
