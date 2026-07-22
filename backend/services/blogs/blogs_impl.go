package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
)

// BlogService はブログに関する処理を提供するサービスインターフェース。
type BlogService interface {
	// FetchBlogs は全ブログデータを取得する。
	FetchBlogs() ([]models.BlogData, error)
	// FetchBlogsByUserId は指定されたユーザーIDに一致するブログデータを取得する。
	FetchBlogsByUserId(userId string) ([]models.BlogData, error)
	// FetchBlogById は指定されたIDに一致するブログデータを取得する。
	FetchBlogById(id string) (*models.BlogData, error)

	// CreateBlog はブログデータを作成する。
	CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// UpdateBlog は指定されたIDに一致するブログデータを更新する。
	UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// DeleteBlog は指定されたIDに一致するブログデータを削除する。
	DeleteBlog(id string) error

	// FetchBlogCategories はブログカテゴリを取得する。
	FetchBlogCategories() ([]string, error)
	// FetchBlogTags はブログタグを取得する。
	FetchBlogTags() ([]string, error)
	// FetchBlogPopular は人気のあるブログを取得する。
	FetchBlogPopular(count int) ([]models.BlogData, error)
}

// BlogServiceImpl は BlogService の実装。
type BlogServiceImpl struct {
	BlogRepository repositories_blogs.BlogRepository
}

// NewBlogService は BlogService インターフェースを実装した BlogServiceImpl を生成する。
func NewBlogService(
	blogRepository repositories_blogs.BlogRepository,
) BlogService {
	return &BlogServiceImpl{
		BlogRepository: blogRepository,
	}
}
