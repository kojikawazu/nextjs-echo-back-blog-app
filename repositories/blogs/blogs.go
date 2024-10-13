package repositories_blogs

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// 全ブログデータを取得する
func (r *BlogRepositoryImpl) FetchBlogs() ([]models.BlogData, error) {
	log.Printf("FetchBlogs start...")

	query := `
        SELECT id, user_id, title, description, github_url, category, tag, likes, created_at, updated_at
        FROM blogs
        ORDER BY created_at DESC
    `

	// Supabaseからクエリを実行し、全データ取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		log.Printf("Failed to fetch blogs: %v", err)
		return nil, err
	}
	log.Println("Fetched blogs successfully")
	defer rows.Close()

	var blogs []models.BlogData

	// 結果をスキャンしてブログデータをリストに追加
	for rows.Next() {
		var blog models.BlogData
		err := rows.Scan(
			&blog.ID,
			&blog.UserId,
			&blog.Title,
			&blog.Description,
			&blog.GithubUrl,
			&blog.Category,
			&blog.Tag,
			&blog.Likes,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan blog: %v", err)
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if rows.Err() != nil {
		log.Printf("Failed to fetch blogs: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Printf("Fetched %d blogs", len(blogs))
	return blogs, nil
}

// 指定されたユーザーIDに一致するブログデータを取得する
func (r *BlogRepositoryImpl) FetchBlogByUserId(userId string) (*models.BlogData, error) {
	log.Printf("FetchBlogByUserId start...")

	query := `
		SELECT id, user_id, title, description, github_url, category, tag, likes, created_at, updated_at
		FROM blogs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	// Supabaseからクエリを実行し、条件に一致するデータを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, userId)

	// 取得した結果をスキャン
	var blog models.BlogData
	err := row.Scan(
		&blog.ID,
		&blog.UserId,
		&blog.Title,
		&blog.Description,
		&blog.GithubUrl,
		&blog.Category,
		&blog.Tag,
		&blog.Likes,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		log.Printf("Blog not found or failed to fetch blog: %v", err)
		return nil, err
	}

	log.Printf("Fetched blog successfully: %v", blog)
	return &blog, nil
}
