package repositories_users

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// 指定されたメールアドレスとパスワードでユーザーを取得する。
// ユーザーが見つからない場合、エラーを返す。
func (r *UserRepositoryImpl) FetchUserByEmailAndPassword(email, password string) (*models.UserData, error) {
	log.Printf("Fetching user from Supabase by email: %s\n", email)

	query := `
        SELECT id, name, email, created_at, updated_at
        FROM users
        WHERE email = $1 AND password = $2
        LIMIT 1
    `

	// Supabaseからクエリを実行し、条件に一致するユーザーを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, email, password)

	// 取得した結果をスキャン
	var user models.UserData
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
