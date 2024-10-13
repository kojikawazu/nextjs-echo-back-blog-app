package repositories_blogs

import "backend/models"

// BlogRepositoryインターフェース
type BlogRepository interface {
	FetchBlogs() ([]models.BlogData, error)
	FetchBlogByUserId(userId string) (*models.BlogData, error)
}

type BlogRepositoryImpl struct{}

// BlogRepositoryインターフェースを実装したBlogRepositoryImplのポインタを返す
func NewBlogRepository() BlogRepository {
	return &BlogRepositoryImpl{}
}
