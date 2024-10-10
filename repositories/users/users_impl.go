package repositories_users

import "backend/models"

// UserRepositoryインターフェース
type UserRepository interface {
	FetchUserByEmailAndPassword(email, password string) (*models.UserData, error)
}

type UserRepositoryImpl struct{}

// UserRepositoryインターフェースを実装したUserRepositoryImplのポインタを返す
func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}
