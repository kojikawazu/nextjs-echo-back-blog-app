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

// 指定されたIDに一致するブログデータを取得する
func (s *BlogServiceImpl) FetchBlogById(id string) (*models.BlogData, error) {
	log.Printf("FetchBlogById start...")

	// バリデーション
	if id == "" {
		log.Printf("invalid id: %s", id)
		return nil, errors.New("invalid id")
	}
	log.Println("Valid id")

	// リポジトリを呼び出してブログデータを取得
	blog, err := s.BlogRepository.FetchBlogById(id)
	if err != nil {
		log.Printf("Failed to fetch blog: %v", err)
		return nil, errors.New("blog not found")
	}

	log.Printf("Fetched blog successfully: %v", blog)
	return blog, nil
}

// ブログデータを作成する
func (s *BlogServiceImpl) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	log.Printf("CreateBlog start...")

	// バリデーション
	if userId == "" {
		log.Printf("invalid userId: %s", userId)
		return nil, errors.New("invalid userId")
	}
	if title == "" {
		log.Printf("invalid title: %s", title)
		return nil, errors.New("invalid title")
	}
	if githubUrl == "" {
		log.Printf("invalid githubUrl: %s", githubUrl)
		return nil, errors.New("invalid githubUrl")
	}
	if category == "" {
		log.Printf("invalid category: %s", category)
		return nil, errors.New("invalid category")
	}
	if description == "" {
		log.Printf("invalid description: %s", description)
		return nil, errors.New("invalid description")
	}
	if tags == "" {
		log.Printf("invalid tags: %s", tags)
		return nil, errors.New("invalid tags")
	}
	log.Println("Valid input")

	// リポジトリを呼び出してブログデータを作成
	blog, err := s.BlogRepository.CreateBlog(userId, title, githubUrl, category, description, tags)
	if err != nil {
		log.Printf("Failed to create blog: %v", err)
		return nil, errors.New("failed to create blog")
	}

	log.Printf("Created blog successfully: %v", blog)
	return blog, nil
}

// 指定されたIDに一致するブログデータを更新する
func (s *BlogServiceImpl) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	log.Printf("UpdateBlog start...")

	// バリデーション
	if id == "" {
		log.Printf("invalid id: %s", id)
		return nil, errors.New("invalid id")
	}
	if title == "" {
		log.Printf("invalid title: %s", title)
		return nil, errors.New("invalid title")
	}
	if githubUrl == "" {
		log.Printf("invalid githubUrl: %s", githubUrl)
		return nil, errors.New("invalid githubUrl")
	}
	if category == "" {
		log.Printf("invalid category: %s", category)
		return nil, errors.New("invalid category")
	}
	if description == "" {
		log.Printf("invalid description: %s", description)
		return nil, errors.New("invalid description")
	}
	if tags == "" {
		log.Printf("invalid tags: %s", tags)
		return nil, errors.New("invalid tags")
	}
	log.Println("Valid input")

	// リポジトリを呼び出してブログデータを更新
	blog, err := s.BlogRepository.UpdateBlog(id, title, githubUrl, category, description, tags)
	if err != nil {
		log.Printf("Failed to update blog: %v", err)
		return nil, errors.New("failed to update blog")
	}

	log.Printf("Updated blog successfully: %v", blog)
	return blog, nil
}

// 指定されたIDに一致するブログデータを削除する
func (s *BlogServiceImpl) DeleteBlog(id string) error {
	log.Printf("DeleteBlog start...")

	// バリデーション
	if id == "" {
		log.Printf("invalid id: %s", id)
		return errors.New("invalid id")
	}
	log.Println("Valid id")

	// リポジトリを呼び出してブログデータを削除
	err := s.BlogRepository.DeleteBlog(id)
	if err != nil {
		log.Printf("Failed to delete blog: %v", err)
		return errors.New("failed to delete blog")
	}

	log.Println("Deleted blog successfully")
	return nil
}

// ブログカテゴリを取得する
func (s *BlogServiceImpl) FetchBlogCategories() ([]string, error) {
	return s.BlogRepository.FetchBlogCategories()
}
