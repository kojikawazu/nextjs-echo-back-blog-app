package services_blog_users

import (
	"backend/models"
	repositories_blog_users "backend/repositories/blog_users"
	"errors"
	"fmt"
	"log"
	"net/mail"
)

// FetchUserByEmailAndPassword は指定されたメールアドレスとパスワードでユーザーを認証し、取得する。
// パスワードは保存されている bcrypt ハッシュと照合する。
// ユーザー不在・パスワード不一致はいずれも ErrInvalidCredentials を返す。
//
// 引数:
//   - email: 取得対象のメールアドレス
//   - password: 取得対象のパスワード
//
// 戻り値:
//   - *models.BlogUsersData: 取得したユーザーデータ
//   - error: 入力が不正な場合は ErrXxx、認証失敗は ErrInvalidCredentials、その他取得に失敗した場合のエラー
func (s *UserServiceImpl) FetchUserByEmailAndPassword(email, password string) (*models.BlogUsersData, error) {
	// バリデーション：emailとpasswordが空でないことを確認
	if email == "" || password == "" {
		log.Printf("Email and password are required")
		return nil, ErrEmailAndPasswordRequired
	}

	// バリデーション：emailが有効な形式であることを確認
	if _, err := mail.ParseAddress(email); err != nil {
		log.Printf("Invalid email format: %v", err)
		return nil, ErrInvalidEmailFormat
	}

	log.Println("Email and password are valid")

	user, err := s.UserRepository.FetchBlogUsersByEmail(email)
	if err != nil {
		if errors.Is(err, repositories_blog_users.ErrUserNotFound) {
			// ユーザー不在でも照合 1 回分の時間をかける（応答時間による登録有無の推測を防ぐ）
			verifyPassword(dummyPasswordHash, password)
			log.Println("User not found")
			return nil, ErrInvalidCredentials
		}
		log.Printf("Failed to fetch user: %v", err)
		return nil, fmt.Errorf("fetch user by email: %w", err)
	}

	if !verifyPassword(user.Password, password) {
		log.Printf("Invalid password: id=%s", user.ID)
		return nil, ErrInvalidCredentials
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
//   - error: id が空の場合は ErrIDRequired、不在の場合は ErrUserNotFound、その他取得に失敗した場合のエラー
func (s *UserServiceImpl) FetchUserById(id string) (*models.BlogUsersData, error) {
	log.Println("Fetching user by id")

	// バリデーション：IDが空でないことを確認
	if id == "" {
		log.Printf("id is required")
		return nil, ErrIDRequired
	}

	log.Println("id is valid")

	user, err := s.UserRepository.FetchBlogUsersById(id)
	if err != nil {
		if errors.Is(err, repositories_blog_users.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("Failed to fetch user: %v", err)
		return nil, fmt.Errorf("fetch user: %w", err)
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
//   - error: 入力が不正な場合は ErrXxx、現在のパスワード不一致は ErrInvalidCurrentPassword、不在の場合は ErrUserNotFound、その他更新に失敗した場合のエラー
func (s *UserServiceImpl) UpdateUser(id, name, email, password, newPassword string) (*models.BlogUsersData, error) {
	log.Println("Updating user")

	// バリデーション：IDが空でないことを確認
	if id == "" {
		log.Printf("id is required")
		return nil, ErrIDRequired
	}
	// バリデーション：nameが空でないことを確認
	if name == "" {
		log.Printf("Name is required")
		return nil, ErrNameRequired
	}
	// バリデーション：emailが空でないことを確認
	if email == "" {
		log.Printf("Email is required")
		return nil, ErrEmailRequired
	}
	// バリデーション：passwordが空でないことを確認
	if password == "" {
		log.Printf("Password is required")
		return nil, ErrPasswordRequired
	}
	// バリデーション：newPasswordが空でないことを確認
	if newPassword == "" {
		log.Printf("New password is required")
		return nil, ErrNewPasswordRequired
	}
	// バリデーション：emailが有効な形式であることを確認
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			log.Printf("Invalid email format: %v", err)
			return nil, ErrInvalidEmailFormat
		}
	}
	// バリデーション：email,passwordを取得し、ユーザーと一致することを確認
	currentUser, err := s.UserRepository.FetchBlogUsersById(id)
	if err != nil {
		if errors.Is(err, repositories_blog_users.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("Failed to validate user: %v", err)
		return nil, fmt.Errorf("fetch user for validation: %w", err)
	}
	if !verifyPassword(currentUser.Password, password) {
		log.Printf("Invalid current password")
		return nil, ErrInvalidCurrentPassword
	}

	log.Println("ID and email are valid")

	// 新しいパスワードは bcrypt でハッシュ化してから保存する（平文を DB に渡さない）
	newPasswordHash, err := hashPassword(newPassword)
	if err != nil {
		if errors.Is(err, ErrPasswordTooLong) {
			log.Printf("New password is too long")
			return nil, ErrPasswordTooLong
		}
		log.Printf("Failed to hash password: %v", err)
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.UserRepository.UpdateBlogUsers(id, name, email, newPasswordHash)
	if err != nil {
		if errors.Is(err, repositories_blog_users.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("Failed to update user: %v", err)
		return nil, fmt.Errorf("update user: %w", err)
	}

	log.Println("Updated user successfully")
	return user, nil
}
