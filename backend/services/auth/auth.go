package services_auth

import (
	"errors"
	"log"
	"net/mail"
)

// ログイン処理を行う
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
