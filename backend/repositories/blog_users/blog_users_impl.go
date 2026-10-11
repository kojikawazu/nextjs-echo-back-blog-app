package repositories_blog_users

import (
	"backend/models"
	"backend/supabase"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
)

// FetchBlogUsersByEmail は指定されたメールアドレスのユーザーを、パスワードハッシュを含めて取得する。
// パスワードの照合は呼び出し側（サービス層）で行う。
//
// 引数:
//   - email: 取得対象のメールアドレス
//
// 戻り値:
//   - *models.BlogUsersData: 取得したユーザー情報（Password は bcrypt ハッシュ）
//   - error: 該当なしの場合は ErrUserNotFound、その他取得に失敗した場合のエラー
func (r *BlogUsersRepositoryImpl) FetchBlogUsersByEmail(email string) (*models.BlogUsersData, error) {
	log.Println("Fetching user from Supabase by email")

	query := `
        SELECT id, name, email, password, created_at, updated_at
        FROM blog_users
        WHERE email = $1
        LIMIT 1
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, email)

	// 取得した結果をスキャン
	var user models.BlogUsersData
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("User not found")
			return nil, ErrUserNotFound
		}
		log.Printf("Failed to fetch user: %v", err)
		return nil, err
	}

	log.Printf("Fetched user successfully: id=%s", user.ID)
	return &user, nil
}

// FetchBlogUsersById は指定されたIDに一致するユーザーを取得する。
//
// 引数:
//   - id: 取得対象のユーザーID
//
// 戻り値:
//   - *models.BlogUsersData: 取得したユーザー情報
//   - error: 取得に失敗した場合のエラー
func (r *BlogUsersRepositoryImpl) FetchBlogUsersById(id string) (*models.BlogUsersData, error) {
	log.Println("Fetching user from Supabase by ID")

	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM blog_users
		WHERE id = $1
		LIMIT 1
	`

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, id)

	// 取得した結果をスキャン
	var user models.BlogUsersData
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("User not found or failed to fetch user: %v", err)
		return nil, err
	}

	log.Printf("Fetched user successfully: id=%s", user.ID)
	return &user, nil
}

// UpdateBlogUsers はユーザー情報を更新する。
//
// 引数:
//   - id: 更新対象のユーザーID
//   - name: 更新後のユーザー名
//   - email: 更新後のメールアドレス
//   - passwordHash: 更新後のパスワードの bcrypt ハッシュ（平文を渡さない）
//
// 戻り値:
//   - *models.BlogUsersData: 更新後のユーザー情報
//   - error: 更新に失敗した場合のエラー
func (r *BlogUsersRepositoryImpl) UpdateBlogUsers(id, name, email, passwordHash string) (*models.BlogUsersData, error) {
	log.Println("Updating user in Supabase")

	query := `
		UPDATE blog_users
		SET name = $1, email = $2, password = $3
		WHERE id = $4
		RETURNING id, name, email, created_at, updated_at
	`

	// Supabaseからクエリを実行し、条件に一致するユーザーを更新
	row := supabase.Pool.QueryRow(supabase.Ctx, query, name, email, passwordHash, id)

	// 取得した結果をスキャン
	var user models.BlogUsersData
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return nil, err
	}

	log.Printf("Updated user successfully: id=%s", user.ID)
	return &user, nil
}
