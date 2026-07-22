package services_auth

// AuthService は認証処理を提供するサービスインターフェース。
type AuthService interface {
	// Login はメール/パスワードを検証してログイン処理を行う。
	//
	// 引数:
	//   - email: 検証対象のメールアドレス
	//   - password: 検証対象のパスワード
	//
	// 戻り値:
	//   - error: バリデーションに失敗した場合のエラー
	Login(email, password string) error
}

// AuthServiceImpl は AuthService の実装。
type AuthServiceImpl struct {
}

// NewAuthService は AuthService インターフェースを実装した AuthServiceImpl を生成する。
//
// 戻り値:
//   - AuthService: 生成された認証サービス
func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}
