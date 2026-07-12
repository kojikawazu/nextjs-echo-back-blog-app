package repositories_auth

// AuthRepositoryインターフェース
type AuthRepository interface{}
type AuthRepositoryImpl struct{}

// AuthRepositoryインターフェースを実装したAuthRepositoryImplのポインタを返す
func NewAuthRepository() AuthRepository {
	return &AuthRepositoryImpl{}
}
