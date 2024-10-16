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

// 指定されたIDに一致するユーザーを取得する
func (s *UserServiceImpl) FetchUserById(id string) (*models.UserData, error) {
	log.Println("Fetching user by id")

	// バリデーション：IDが空でないことを確認
	if id == "" {
		log.Printf("id is required")
		return nil, errors.New("id is required")
	}

	log.Println("id is valid")

	user, err := s.UserRepository.FetchUserById(id)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, errors.New("failed to fetch user")
	}

	log.Println("Fetched user successfully")
	return user, nil
}

// 指定されたIDに一致するユーザーを更新する
func (s *UserServiceImpl) UpdateUser(id, name, email, password, newPassword string) (*models.UserData, error) {
	log.Println("Updating user")

	// バリデーション：IDが空でないことを確認
	if id == "" {
		log.Printf("id is required")
		return nil, errors.New("id is required")
	}
	// バリデーション：nameが空でないことを確認
	if name == "" {
		log.Printf("Name is required")
		return nil, errors.New("name is required")
	}
	// バリデーション：emailが空でないことを確認
	if email == "" {
		log.Printf("Email is required")
		return nil, errors.New("email is required")
	}
	// バリデーション：passwordが空でないことを確認
	if password == "" {
		log.Printf("Password is required")
		return nil, errors.New("password is required")
	}
	// バリデーション：newPasswordが空でないことを確認
	if newPassword == "" {
		log.Printf("New password is required")
		return nil, errors.New("new password is required")
	}
	// バリデーション：emailが有効な形式であることを確認
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			log.Printf("Invalid email format: %v", err)
			return nil, errors.New("invalid email format")
		}
	}
	// バリデーション：email,passwordを取得し、ユーザーと一致することを確認
	currentUser, err := s.UserRepository.FetchUserById(id)
	if err != nil {
		log.Printf("Failed to validate user: %v", err)
		return nil, errors.New("failed to validate user")
	}
	if currentUser.Password != password {
		log.Printf("Invalid current password")
		return nil, errors.New("invalid current password")
	}

	log.Println("ID and email are valid")

	user, err := s.UserRepository.UpdateUser(id, name, email, newPassword)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, errors.New("failed to update user")
	}

	log.Println("Updated user successfully")
	return user, nil
}
