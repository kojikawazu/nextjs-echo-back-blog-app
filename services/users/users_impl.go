package services_users

import (
	"backend/models"
	repositories_users "backend/repositories/users"
)

// UserServiceインターフェース
type UserService interface {
	FetchUserByEmailAndPassword(email, password string) (*models.UserData, error)
	FetchUserById(id string) (*models.UserData, error)
	UpdateUser(id, name, email, password string) (*models.UserData, error)
}
type UserServiceImpl struct {
	UserRepository repositories_users.UserRepository
}

// UserServiceインターフェースを実装したUserServiceImplのポインタを返す
func NewUserService(
	userRepository repositories_users.UserRepository,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}
