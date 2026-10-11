package repositories_blog_users

import (
	"backend/models"
	"errors"
)

// ErrUserNotFound は該当するユーザーが存在しない場合のエラー。
// DB ドライバ固有のエラー（pgx.ErrNoRows）を呼び出し側に漏らさないため、リポジトリ層で変換する。
var ErrUserNotFound = errors.New("user not found")

// BlogUsersRepository はユーザーのデータアクセスを担うリポジトリインターフェース。
type BlogUsersRepository interface {
	// FetchBlogUsersByEmail は指定されたメールアドレスのユーザーを、パスワードハッシュを含めて取得する。
	// パスワードの照合は呼び出し側（サービス層）で行う。
	//
	// 引数:
	//   - email: 取得対象のメールアドレス
	//
	// 戻り値:
	//   - *models.BlogUsersData: 取得したユーザー情報（Password は bcrypt ハッシュ）
	//   - error: 該当なしの場合は ErrUserNotFound、その他取得に失敗した場合のエラー
	FetchBlogUsersByEmail(email string) (*models.BlogUsersData, error)
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
	//   - passwordHash: 更新後のパスワードの bcrypt ハッシュ（平文を渡さない）
	//
	// 戻り値:
	//   - *models.BlogUsersData: 更新後のユーザー情報
	//   - error: 更新に失敗した場合のエラー
	UpdateBlogUsers(id, name, email, passwordHash string) (*models.BlogUsersData, error)
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
