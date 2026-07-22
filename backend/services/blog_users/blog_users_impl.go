package services_blog_users

import (
	"backend/models"
	repositories_blog_users "backend/repositories/blog_users"
)

// UserService はユーザーに関する処理を提供するサービスインターフェース。
type UserService interface {
	// FetchUserByEmailAndPassword は指定されたメールアドレスとパスワードでユーザーを取得する。
	FetchUserByEmailAndPassword(email, password string) (*models.BlogUsersData, error)
	// FetchUserById は指定されたIDに一致するユーザーを取得する。
	FetchUserById(id string) (*models.BlogUsersData, error)
	// UpdateUser は指定されたIDに一致するユーザーを更新する。
	UpdateUser(id, name, email, password, newPassword string) (*models.BlogUsersData, error)
}

// UserServiceImpl は UserService の実装。
type UserServiceImpl struct {
	UserRepository repositories_blog_users.BlogUsersRepository
}

// NewUserService は UserService インターフェースを実装した UserServiceImpl を生成する。
func NewUserService(
	userRepository repositories_blog_users.BlogUsersRepository,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}
