package repositories_blog_users

import "backend/models"

// BlogUsersRepositoryインターフェース
type BlogUsersRepository interface {
	FetchBlogUsersByEmailAndPassword(email, password string) (*models.BlogUsersData, error)
	FetchBlogUsersById(id string) (*models.BlogUsersData, error)
	UpdateBlogUsers(id, name, email, password string) (*models.BlogUsersData, error)
}

type BlogUsersRepositoryImpl struct{}

// BlogUsersRepositoryインターフェースを実装したBlogUsersRepositoryImplのポインタを返す
func NewBlogUsersRepository() BlogUsersRepository {
	return &BlogUsersRepositoryImpl{}
}
