package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
)

// BlogServiceインターフェース
type BlogService interface {
	FetchBlogs() ([]models.BlogData, error)
	FetchBlogsByUserId(userId string) ([]models.BlogData, error)
	FetchBlogById(id string) (*models.BlogData, error)
	CreateBlog(userId, title, githubUrl, category, description, tags string) (models.BlogData, error)
	UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error)
	DeleteBlog(id string) error
}

type BlogServiceImpl struct {
	BlogRepository repositories_blogs.BlogRepository
}

// BlogServiceインターフェースを実装したBlogServiceImplのポインタを返す
func NewBlogService(
	blogRepository repositories_blogs.BlogRepository,
) BlogService {
	return &BlogServiceImpl{
		BlogRepository: blogRepository,
	}
}
