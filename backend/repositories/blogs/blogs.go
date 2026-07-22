package repositories_blogs

import "backend/models"

// BlogRepository はブログのデータアクセスを担うリポジトリインターフェース。
type BlogRepository interface {
	// FetchBlogs は全ブログデータを取得する。
	FetchBlogs() ([]models.BlogData, error)
	// FetchBlogsByUserId は指定されたユーザーIDに一致するブログデータを取得する。
	FetchBlogsByUserId(userId string) ([]models.BlogData, error)
	// FetchBlogById は指定されたIDに一致するブログデータを取得する。
	FetchBlogById(id string) (*models.BlogData, error)

	// CreateBlog はブログデータを作成する。
	CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// UpdateBlog はブログデータを更新する。
	UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	// DeleteBlog はブログデータを削除する。
	DeleteBlog(id string) error

	// FetchBlogCategories はブログカテゴリ一覧を取得する。
	FetchBlogCategories() ([]string, error)
	// FetchBlogTags はブログタグ一覧を取得する。
	FetchBlogTags() ([]string, error)
	// FetchBlogPopular は人気のあるブログを取得する。
	FetchBlogPopular(count int) ([]models.BlogData, error)
}

// BlogRepositoryImpl は BlogRepository の実装。
type BlogRepositoryImpl struct{}

// NewBlogRepository は BlogRepository を実装した BlogRepositoryImpl のポインタを生成する。
func NewBlogRepository() BlogRepository {
	return &BlogRepositoryImpl{}
}
