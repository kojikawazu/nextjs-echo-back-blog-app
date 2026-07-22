package services_blog_users

import (
	"backend/models"
	"database/sql"
	"errors"
	"log"
	"net/mail"
)

// FetchUserByEmailAndPassword は指定されたメールアドレスとパスワードでユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
//
// 引数:
//   - email: 取得対象のメールアドレス
//   - password: 取得対象のパスワード
//
// 戻り値:
//   - *models.BlogUsersData: 取得したユーザーデータ
//   - error: 取得に失敗した場合のエラー
func (s *UserServiceImpl) FetchUserByEmailAndPassword(email, password string) (*models.BlogUsersData, error) {
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

	user, err := s.UserRepository.FetchBlogUsersByEmailAndPassword(email, password)
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

// FetchUserById は指定されたIDに一致するユーザーを取得する。
//
// 引数:
//   - id: 取得対象のユーザーID
//
// 戻り値:
//   - *models.BlogUsersData: 取得したユーザーデータ
//   - error: 取得に失敗した場合のエラー
func (s *UserServiceImpl) FetchUserById(id string) (*models.BlogUsersData, error) {
	log.Println("Fetching user by id")

	// バリデーション：IDが空でないことを確認
	if id == "" {
		log.Printf("id is required")
		return nil, errors.New("id is required")
	}

	log.Println("id is valid")

	user, err := s.UserRepository.FetchBlogUsersById(id)
	if err != nil {
		log.Printf("Failed to fetch user: %v", err)
		return nil, errors.New("failed to fetch user")
	}

	log.Println("Fetched user successfully")
	return user, nil
}

// UpdateUser は指定されたIDに一致するユーザーを更新する。
//
// 引数:
//   - id: 更新対象のユーザーID
//   - name: 更新後のユーザー名
//   - email: 更新後のメールアドレス
//   - password: 本人確認用の現在のパスワード
//   - newPassword: 更新後の新しいパスワード
//
// 戻り値:
//   - *models.BlogUsersData: 更新後のユーザーデータ
//   - error: 更新に失敗した場合のエラー
func (s *UserServiceImpl) UpdateUser(id, name, email, password, newPassword string) (*models.BlogUsersData, error) {
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
	currentUser, err := s.UserRepository.FetchBlogUsersById(id)
	if err != nil {
		log.Printf("Failed to validate user: %v", err)
		return nil, errors.New("failed to validate user")
	}
	if currentUser.Password != password {
		log.Printf("Invalid current password")
		return nil, errors.New("invalid current password")
	}

	log.Println("ID and email are valid")

	user, err := s.UserRepository.UpdateBlogUsers(id, name, email, newPassword)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, errors.New("failed to update user")
	}

	log.Println("Updated user successfully")
	return user, nil
}
