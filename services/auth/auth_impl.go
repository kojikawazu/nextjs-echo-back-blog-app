package services_auth

// AuthServiceインターフェース
type AuthService interface {
	Login(email, password string) error
}
type AuthServiceImpl struct {
}

// AuthServiceインターフェースを実装したAuthServiceImplのポインタを返す
func NewAuthService() AuthService {
	return &AuthServiceImpl{}
}
