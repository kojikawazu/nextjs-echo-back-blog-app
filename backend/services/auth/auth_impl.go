package services_auth

// AuthService は認証処理を提供するサービスインターフェース。
type AuthService interface {
	// Login はメール/パスワードを検証してログイン処理を行う。
	Login(email, password string) error
}

// AuthServiceImpl は AuthService の実装。
type AuthServiceImpl struct {
}

// NewAuthService は AuthService インターフェースを実装した AuthServiceImpl を生成する。
func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}
