package services_blogs

import (
	"backend/logger"
	"backend/models"
	"errors"
)

// 全ブログデータを取得する
func (s *BlogServiceImpl) FetchBlogs() ([]models.BlogData, error) {
	return s.BlogRepository.FetchBlogs()
}

// 指定されたユーザーIDに一致するブログデータを取得する
func (s *BlogServiceImpl) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogsByUserId start...")

	// バリデーション
	if userId == "" {
		logger.ErrorLog.Printf("invalid userId: %s", userId)
		return nil, errors.New("invalid userId")
	}
	logger.InfoLog.Println("Valid userId")

	// リポジトリを呼び出してブログデータを取得
	blogs, err := s.BlogRepository.FetchBlogsByUserId(userId)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blogs: %v", err)
		return nil, errors.New("blogs not found")
	}

	logger.InfoLog.Printf("Fetched blogs successfully: %v", blogs)
	return blogs, nil
}

// 指定されたIDに一致するブログデータを取得する
func (s *BlogServiceImpl) FetchBlogById(id string) (*models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogById start...")

	// バリデーション
	if id == "" {
		logger.ErrorLog.Printf("invalid id: %s", id)
		return nil, errors.New("invalid id")
	}
	logger.InfoLog.Println("Valid id")

	// リポジトリを呼び出してブログデータを取得
	blog, err := s.BlogRepository.FetchBlogById(id)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blog: %v", err)
		return nil, errors.New("blog not found")
	}

	logger.InfoLog.Printf("Fetched blog successfully: %v", blog)
	return blog, nil
}

// ブログデータを作成する
func (s *BlogServiceImpl) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	logger.InfoLog.Printf("CreateBlog start...")

	// バリデーション
	if userId == "" {
		logger.ErrorLog.Printf("invalid userId: %s", userId)
		return nil, errors.New("invalid userId")
	}
	if title == "" {
		logger.ErrorLog.Printf("invalid title: %s", title)
		return nil, errors.New("invalid title")
	}
	if githubUrl == "" {
		logger.ErrorLog.Printf("invalid githubUrl: %s", githubUrl)
		return nil, errors.New("invalid githubUrl")
	}
	if category == "" {
		logger.ErrorLog.Printf("invalid category: %s", category)
		return nil, errors.New("invalid category")
	}
	if description == "" {
		logger.ErrorLog.Printf("invalid description: %s", description)
		return nil, errors.New("invalid description")
	}
	if tags == "" {
		logger.ErrorLog.Printf("invalid tags: %s", tags)
		return nil, errors.New("invalid tags")
	}
	logger.InfoLog.Println("Valid input")

	// リポジトリを呼び出してブログデータを作成
	blog, err := s.BlogRepository.CreateBlog(userId, title, githubUrl, category, description, tags)
	if err != nil {
		logger.ErrorLog.Printf("Failed to create blog: %v", err)
		return nil, errors.New("failed to create blog")
	}

	logger.InfoLog.Printf("Created blog successfully: %v", blog)
	return blog, nil
}

// 指定されたIDに一致するブログデータを更新する
func (s *BlogServiceImpl) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	logger.InfoLog.Printf("UpdateBlog start...")

	// バリデーション
	if id == "" {
		logger.ErrorLog.Printf("invalid id: %s", id)
		return nil, errors.New("invalid id")
	}
	if title == "" {
		logger.ErrorLog.Printf("invalid title: %s", title)
		return nil, errors.New("invalid title")
	}
	if githubUrl == "" {
		logger.ErrorLog.Printf("invalid githubUrl: %s", githubUrl)
		return nil, errors.New("invalid githubUrl")
	}
	if category == "" {
		logger.ErrorLog.Printf("invalid category: %s", category)
		return nil, errors.New("invalid category")
	}
	if description == "" {
		logger.ErrorLog.Printf("invalid description: %s", description)
		return nil, errors.New("invalid description")
	}
	if tags == "" {
		logger.ErrorLog.Printf("invalid tags: %s", tags)
		return nil, errors.New("invalid tags")
	}
	logger.InfoLog.Println("Valid input")

	// リポジトリを呼び出してブログデータを更新
	blog, err := s.BlogRepository.UpdateBlog(id, title, githubUrl, category, description, tags)
	if err != nil {
		logger.ErrorLog.Printf("Failed to update blog: %v", err)
		return nil, errors.New("failed to update blog")
	}

	logger.InfoLog.Printf("Updated blog successfully: %v", blog)
	return blog, nil
}

// 指定されたIDに一致するブログデータを削除する
func (s *BlogServiceImpl) DeleteBlog(id string) error {
	logger.InfoLog.Printf("DeleteBlog start...")

	// バリデーション
	if id == "" {
		logger.ErrorLog.Printf("invalid id: %s", id)
		return errors.New("invalid id")
	}
	logger.InfoLog.Println("Valid id")

	// リポジトリを呼び出してブログデータを削除
	err := s.BlogRepository.DeleteBlog(id)
	if err != nil {
		logger.ErrorLog.Printf("Failed to delete blog: %v", err)
		return errors.New("failed to delete blog")
	}

	logger.InfoLog.Println("Deleted blog successfully")
	return nil
}

// ブログカテゴリを取得する
func (s *BlogServiceImpl) FetchBlogCategories() ([]string, error) {
	return s.BlogRepository.FetchBlogCategories()
}

// ブログタグを取得する
func (s *BlogServiceImpl) FetchBlogTags() ([]string, error) {
	return s.BlogRepository.FetchBlogTags()
}

// 人気のあるブログを取得する
func (s *BlogServiceImpl) FetchBlogPopular(count int) ([]models.BlogData, error) {

	// バリデーション
	if count <= 0 {
		logger.ErrorLog.Printf("invalid count: %d", count)
		return nil, errors.New("invalid count")
	}
	logger.InfoLog.Println("Valid count")

	// リポジトリを呼び出して人気のあるブログを取得
	blogs, err := s.BlogRepository.FetchBlogPopular(count)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch popular blogs: %v", err)
		return nil, errors.New("failed to fetch popular blogs")
	}

	logger.InfoLog.Printf("Fetched popular blogs successfully: %v", blogs)
	return blogs, nil
}
