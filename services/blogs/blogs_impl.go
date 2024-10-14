package services_blogs

import (
	"backend/models"
	repositories_blogs "backend/repositories/blogs"
)

// BlogServiceインターフェース
type BlogService interface {
	FetchBlogs() ([]models.BlogData, error)
	FetchBlogsByUserId(userId string) ([]models.BlogData, error)
	CreateBlog(userId, title, githubUrl, category, description, tags string) (models.BlogData, error)
}

type BlogServiceImpl struct {
	BlogRepository repositories_blogs.BlogRepository
}

// UserServiceインターフェースを実装したUserServiceImplのポインタを返す
func NewBlogService(
	blogRepository repositories_blogs.BlogRepository,
) BlogService {
	return &BlogServiceImpl{
		BlogRepository: blogRepository,
	}
}
