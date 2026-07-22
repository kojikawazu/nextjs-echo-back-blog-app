package repositories_auth

// AuthRepository は認証系のデータアクセスを担うリポジトリインターフェース。
type AuthRepository interface{}

// AuthRepositoryImpl は AuthRepository の実装。
type AuthRepositoryImpl struct{}

// NewAuthRepository は AuthRepository を実装した AuthRepositoryImpl のポインタを生成する。
//
// 戻り値:
//   - AuthRepository: 生成したリポジトリ実装
func NewAuthRepository() AuthRepository {
	return &AuthRepositoryImpl{}
}
