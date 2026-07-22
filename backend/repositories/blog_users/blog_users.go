package repositories_blog_users

import "backend/models"

// BlogUsersRepository はユーザーのデータアクセスを担うリポジトリインターフェース。
type BlogUsersRepository interface {
	// FetchBlogUsersByEmailAndPassword は指定されたメールアドレスとパスワードでユーザーを取得する。
	//
	// 引数:
	//   - email: 取得対象のメールアドレス
	//   - password: 認証に用いるパスワード
	//
	// 戻り値:
	//   - *models.BlogUsersData: 取得したユーザー情報
	//   - error: 取得に失敗した場合のエラー
	FetchBlogUsersByEmailAndPassword(email, password string) (*models.BlogUsersData, error)
	// FetchBlogUsersById は指定されたIDに一致するユーザーを取得する。
	//
	// 引数:
	//   - id: 取得対象のユーザーID
	//
	// 戻り値:
	//   - *models.BlogUsersData: 取得したユーザー情報
	//   - error: 取得に失敗した場合のエラー
	FetchBlogUsersById(id string) (*models.BlogUsersData, error)
	// UpdateBlogUsers はユーザー情報を更新する。
	//
	// 引数:
	//   - id: 更新対象のユーザーID
	//   - name: 更新後のユーザー名
	//   - email: 更新後のメールアドレス
	//   - password: 更新後のパスワード
	//
	// 戻り値:
	//   - *models.BlogUsersData: 更新後のユーザー情報
	//   - error: 更新に失敗した場合のエラー
	UpdateBlogUsers(id, name, email, password string) (*models.BlogUsersData, error)
}

// BlogUsersRepositoryImpl は BlogUsersRepository の実装。
type BlogUsersRepositoryImpl struct{}

// NewBlogUsersRepository は BlogUsersRepository を実装した BlogUsersRepositoryImpl のポインタを生成する。
//
// 戻り値:
//   - BlogUsersRepository: 生成したリポジトリ実装
func NewBlogUsersRepository() BlogUsersRepository {
	return &BlogUsersRepositoryImpl{}
}
