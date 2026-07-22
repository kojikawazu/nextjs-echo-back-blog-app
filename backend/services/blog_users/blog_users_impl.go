package services_blog_users

import (
	"backend/models"
	repositories_blog_users "backend/repositories/blog_users"
)

// UserService はユーザーに関する処理を提供するサービスインターフェース。
type UserService interface {
	// FetchUserByEmailAndPassword は指定されたメールアドレスとパスワードでユーザーを取得する。
	//
	// 引数:
	//   - email: 取得対象のメールアドレス
	//   - password: 取得対象のパスワード
	//
	// 戻り値:
	//   - *models.BlogUsersData: 取得したユーザーデータ
	//   - error: 取得に失敗した場合のエラー
	FetchUserByEmailAndPassword(email, password string) (*models.BlogUsersData, error)
	// FetchUserById は指定されたIDに一致するユーザーを取得する。
	//
	// 引数:
	//   - id: 取得対象のユーザーID
	//
	// 戻り値:
	//   - *models.BlogUsersData: 取得したユーザーデータ
	//   - error: 取得に失敗した場合のエラー
	FetchUserById(id string) (*models.BlogUsersData, error)
	// UpdateUser は指定されたIDに一致するユーザーを更新する。
	//
	// 引数:
	//   - id: 更新対象のユーザーID
	//   - name: 更新後のユーザー名
	//   - email: 更新後のメールアドレス
	//   - password: 本人確認用の現在のパスワード
	//   - newPassword: 更新後の新しいパスワード
	//
	// 戻り値:
	//   - *models.BlogUsersData: 更新後のユーザーデータ
	//   - error: 更新に失敗した場合のエラー
	UpdateUser(id, name, email, password, newPassword string) (*models.BlogUsersData, error)
}

// UserServiceImpl は UserService の実装。
type UserServiceImpl struct {
	UserRepository repositories_blog_users.BlogUsersRepository
}

// NewUserService は UserService インターフェースを実装した UserServiceImpl を生成する。
//
// 引数:
//   - userRepository: ユーザーデータへのアクセスを担うリポジトリ
//
// 戻り値:
//   - UserService: 生成されたユーザーサービス
func NewUserService(
	userRepository repositories_blog_users.BlogUsersRepository,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}
