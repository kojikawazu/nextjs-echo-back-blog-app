package services_blogs

import (
	"backend/models"
	"errors"
	"log"
)

// 全ブログデータを取得する
func (s *BlogServiceImpl) FetchBlogs() ([]models.BlogData, error) {
	return s.BlogRepository.FetchBlogs()
}

// 指定されたユーザーIDに一致するブログデータを取得する
func (s *BlogServiceImpl) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	log.Printf("FetchBlogsByUserId start...")

	// バリデーション
	if userId == "" {
		log.Printf("invalid userId: %s", userId)
		return nil, errors.New("invalid userId")
	}
	log.Println("Valid userId")

	// リポジトリを呼び出してブログデータを取得
	blogs, err := s.BlogRepository.FetchBlogsByUserId(userId)
	if err != nil {
		log.Printf("Failed to fetch blogs: %v", err)
		return nil, errors.New("blogs not found")
	}

	log.Printf("Fetched blogs successfully: %v", blogs)
	return blogs, nil
}
