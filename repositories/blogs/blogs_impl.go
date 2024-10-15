package repositories_blogs

import "backend/models"

// BlogRepositoryインターフェース
type BlogRepository interface {
	FetchBlogs() ([]models.BlogData, error)
	FetchBlogsByUserId(userId string) ([]models.BlogData, error)
	FetchBlogById(id string) (*models.BlogData, error)
	CreateBlog(userId, title, githubUrl, category, description, tags string) (models.BlogData, error)
	UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	DeleteBlog(id string) error
}

type BlogRepositoryImpl struct{}

// BlogRepositoryインターフェースを実装したBlogRepositoryImplのポインタを返す
func NewBlogRepository() BlogRepository {
	return &BlogRepositoryImpl{}
}
