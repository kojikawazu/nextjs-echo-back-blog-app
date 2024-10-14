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

// ブログデータを作成する
func (s *BlogServiceImpl) CreateBlog(userId, title, githubUrl, category, description, tags string) (models.BlogData, error) {
	log.Printf("CreateBlog start...")

	// バリデーション
	if userId == "" {
		log.Printf("invalid userId: %s", userId)
		return models.BlogData{}, errors.New("invalid userId")
	}
	if title == "" {
		log.Printf("invalid title: %s", title)
		return models.BlogData{}, errors.New("invalid title")
	}
	if githubUrl == "" {
		log.Printf("invalid githubUrl: %s", githubUrl)
		return models.BlogData{}, errors.New("invalid githubUrl")
	}
	if category == "" {
		log.Printf("invalid category: %s", category)
		return models.BlogData{}, errors.New("invalid category")
	}
	if description == "" {
		log.Printf("invalid description: %s", description)
		return models.BlogData{}, errors.New("invalid description")
	}
	if tags == "" {
		log.Printf("invalid tags: %s", tags)
		return models.BlogData{}, errors.New("invalid tags")
	}
	log.Println("Valid input")

	// リポジトリを呼び出してブログデータを作成
	blog, err := s.BlogRepository.CreateBlog(userId, title, githubUrl, category, description, tags)
	if err != nil {
		log.Printf("Failed to create blog: %v", err)
		return models.BlogData{}, errors.New("failed to create blog")
	}

	log.Printf("Created blog successfully: %v", blog)
	return blog, nil
}
