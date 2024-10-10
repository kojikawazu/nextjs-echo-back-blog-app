package services_users

import (
	"backend/models"
	"database/sql"
	"errors"
	"log"
	"net/mail"
)

// 指定されたメールアドレスとパスワードでユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (s *UserServiceImpl) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	// バリデーション：emailとpasswordが空でないことを確認
	if email == "" || password == "" {
		log.Printf("Email and password are required")
		return nil, errors.New("email and password are required")
	}

	// バリデーション：emailが有効な形式であることを確認
	if _, err := mail.ParseAddress(email); err != nil {
		log.Printf("Invalid email format: %v", err)
		return nil, errors.New("invalid email format")
	}

	log.Println("Email and password are valid")

	user, err := s.UserRepository.FetchUserByEmailAndPassword(email, password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found for email: %s", email)
			return nil, errors.New("user not found")
		}
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}

	return user, nil
}
