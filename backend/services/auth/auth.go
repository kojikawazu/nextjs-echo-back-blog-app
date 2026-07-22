package services_auth

import (
	"errors"
	"log"
	"net/mail"
)

// Login はメール/パスワードを検証してログイン処理を行う。
//
// 引数:
//   - email: 検証対象のメールアドレス
//   - password: 検証対象のパスワード
//
// 戻り値:
//   - error: バリデーションに失敗した場合のエラー
func (r *AuthServiceImpl) Login(email, password string) error {
	log.Println("Logging in...")

	// バリデーション：emailとpasswordが空でないことを確認
	if email == "" || password == "" {
		log.Println("Email and password are required")
		return errors.New("email and password are required")
	}
	// バリデーション：emailが有効な形式であることを確認
	if _, err := mail.ParseAddress(email); err != nil {
		log.Println("Invalid email format")
		return errors.New("invalid email format")
	}

	log.Println("Email and password are valid")
	return nil
}
