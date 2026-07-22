package repositories_blog_users

import "backend/models"

// BlogUsersRepository はユーザーのデータアクセスを担うリポジトリインターフェース。
type BlogUsersRepository interface {
	// FetchBlogUsersByEmailAndPassword は指定されたメールアドレスとパスワードでユーザーを取得する。
	FetchBlogUsersByEmailAndPassword(email, password string) (*models.BlogUsersData, error)
	// FetchBlogUsersById は指定されたIDに一致するユーザーを取得する。
	FetchBlogUsersById(id string) (*models.BlogUsersData, error)
	// UpdateBlogUsers はユーザー情報を更新する。
	UpdateBlogUsers(id, name, email, password string) (*models.BlogUsersData, error)
}

// BlogUsersRepositoryImpl は BlogUsersRepository の実装。
type BlogUsersRepositoryImpl struct{}

// NewBlogUsersRepository は BlogUsersRepository を実装した BlogUsersRepositoryImpl のポインタを生成する。
func NewBlogUsersRepository() BlogUsersRepository {
	return &BlogUsersRepositoryImpl{}
}
