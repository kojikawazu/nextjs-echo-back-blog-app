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
        SELECT id, user_id, title, description, github_url, category, tags, likes, created_at, updated_at
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
			&blog.Tags,
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
func (r *BlogRepositoryImpl) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	log.Printf("FetchBlogsByUserId start...")

	query := `
		SELECT id, user_id, title, description, github_url, category, tags, likes, created_at, updated_at
		FROM blogs
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	// Supabaseからクエリを実行し、条件に一致するデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query, userId)
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
			&blog.Tags,
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

// 指定されたIDに一致するブログデータを取得する
func (r *BlogRepositoryImpl) FetchBlogById(id string) (*models.BlogData, error) {
	log.Printf("FetchBlogById start...")

	query := `
		SELECT id, user_id, title, description, github_url, category, tags, likes, created_at, updated_at
		FROM blogs
		WHERE id = $1
	`

	// Supabaseからクエリを実行し、条件に一致するデータを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, id)
	var blog models.BlogData
	err := row.Scan(
		&blog.ID,
		&blog.UserId,
		&blog.Title,
		&blog.Description,
		&blog.GithubUrl,
		&blog.Category,
		&blog.Tags,
		&blog.Likes,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to fetch blog: %v", err)
		return nil, err
	}

	log.Printf("Fetched blog: %v", blog)
	return &blog, nil
}

// ブログデータの作成
func (r *BlogRepositoryImpl) CreateBlog(userId, title, githubUrl, category, description, tags string) (models.BlogData, error) {
	log.Printf("CreateReservation start...")

	query := `
		INSERT INTO blogs (user_id, title, github_url, category, description, tags)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, description, github_url, category, tags, likes, created_at, updated_at
	`

	// Supabaseからクエリを実行し、新しいブログデータを作成
	row := supabase.Pool.QueryRow(supabase.Ctx, query, userId, title, githubUrl, category, description, tags)
	// 結果をスキャンして新しいブログデータを返す
	var blog models.BlogData
	err := row.Scan(
		&blog.ID,
		&blog.UserId,
		&blog.Title,
		&blog.Description,
		&blog.GithubUrl,
		&blog.Category,
		&blog.Tags,
		&blog.Likes,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to create blog: %v", err)
		return models.BlogData{}, err
	}

	log.Printf("Created blog: %v", blog)
	return blog, nil
}
