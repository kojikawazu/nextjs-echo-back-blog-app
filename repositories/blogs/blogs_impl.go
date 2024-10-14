package repositories_blogs

import "backend/models"

// BlogRepositoryインターフェース
type BlogRepository interface {
	FetchBlogs() ([]models.BlogData, error)
	FetchBlogsByUserId(userId string) ([]models.BlogData, error)
	CreateBlog(userId, title, githubUrl, category, description, tags string) (models.BlogData, error)
}

type BlogRepositoryImpl struct{}

// BlogRepositoryインターフェースを実装したBlogRepositoryImplのポインタを返す
func NewBlogRepository() BlogRepository {
	return &BlogRepositoryImpl{}
}
