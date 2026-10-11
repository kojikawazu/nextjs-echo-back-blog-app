package services_auth

import (
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
//   - error: 未入力の場合は ErrEmailAndPasswordRequired、メールアドレスの形式が不正な場合は ErrInvalidEmailFormat
func (r *AuthServiceImpl) Login(email, password string) error {
	log.Println("Logging in...")

	// バリデーション：emailとpasswordが空でないことを確認
	if email == "" || password == "" {
		log.Println("Email and password are required")
		return ErrEmailAndPasswordRequired
	}
	// バリデーション：emailが有効な形式であることを確認
	if _, err := mail.ParseAddress(email); err != nil {
		log.Println("Invalid email format")
		return ErrInvalidEmailFormat
	}

	log.Println("Email and password are valid")
	return nil
}
