package repositories_blog_users

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// FetchBlogUsersByEmailAndPassword は指定されたメールアドレスとパスワードでユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
//
// 引数:
//   - email: 取得対象のメールアドレス
//   - password: 認証に用いるパスワード
//
// 戻り値:
//   - *models.BlogUsersData: 取得したユーザー情報
//   - error: 取得に失敗した場合のエラー
func (r *BlogUsersRepositoryImpl) FetchBlogUsersByEmailAndPassword(email, password string) (*models.BlogUsersData, error) {
	log.Printf("Fetching user from Supabase by email: %s\n", email)

	query := `
        SELECT id, name, email, created_at, updated_at
        FROM blog_users
        WHERE email = $1 AND password = $2
        LIMIT 1
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, email, password)

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
		log.Printf("User not found or failed to fetch user: %v", err)
		return nil, err
	}

	log.Printf("Fetched user successfully: %v", user)
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

	log.Printf("Fetched user successfully: %v", user)
	return &user, nil
}

// UpdateBlogUsers はユーザー情報を更新する。
//
// 引数:
//   - id: 更新対象のユーザーID
//   - name: 更新後のユーザー名
//   - email: 更新後のメールアドレス
//   - password: 更新後のパスワード
//
// 戻り値:
//   - *models.BlogUsersData: 更新後のユーザー情報
//   - error: 更新に失敗した場合のエラー
func (r *BlogUsersRepositoryImpl) UpdateBlogUsers(id, name, email, password string) (*models.BlogUsersData, error) {
	log.Println("Updating user in Supabase")

	query := `
		UPDATE blog_users
		SET name = $1, email = $2, password = $3
		WHERE id = $4
		RETURNING id, name, email, created_at, updated_at
	`

	// Supabaseからクエリを実行し、条件に一致するユーザーを更新
	row := supabase.Pool.QueryRow(supabase.Ctx, query, name, email, password, id)

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

	log.Printf("Updated user successfully: %v", user)
	return &user, nil
}
