package services_blog_users

import (
	"backend/models"
	repositories_blog_users "backend/repositories/blog_users"
)

// UserServiceインターフェース
type UserService interface {
	FetchUserByEmailAndPassword(email, password string) (*models.BlogUsersData, error)
	FetchUserById(id string) (*models.BlogUsersData, error)
	UpdateUser(id, name, email, password, newPassword string) (*models.BlogUsersData, error)
}
type UserServiceImpl struct {
	UserRepository repositories_blog_users.BlogUsersRepository
}

// UserServiceインターフェースを実装したUserServiceImplのポインタを返す
func NewUserService(
	userRepository repositories_blog_users.BlogUsersRepository,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}
