package repositories_blogs

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// 全ブログデータを取得する
func (r *BlogRepositoryImpl) FetchBlogs() ([]models.BlogData, error) {
	log.Printf("FetchBlogs start...")

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
        SELECT b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags,
               COALESCE(COUNT(bl.id), 0) AS likes,
			   b.created_at, b.updated_at
        FROM blogs b
        LEFT JOIN blogs_likes bl ON b.id = bl.blog_id
        GROUP BY b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags, b.created_at, b.updated_at
        ORDER BY b.created_at DESC
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
		var likeCount int

		err := rows.Scan(
			&blog.ID,
			&blog.UserId,
			&blog.Title,
			&blog.Description,
			&blog.GithubUrl,
			&blog.Category,
			&blog.Tags,
			&likeCount,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan blog: %v", err)
			return nil, err
		}

		// `Likes` フィールドにキャストして代入
		blog.Likes = int8(likeCount)

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

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
		SELECT b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags, 
		       COALESCE(COUNT(bl.id), 0) AS likes, b.created_at, b.updated_at
		FROM blogs b
		LEFT JOIN blogs_likes bl ON b.id = bl.blog_id
		WHERE b.user_id = $1
		GROUP BY b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags, b.created_at, b.updated_at
		ORDER BY b.created_at DESC
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
		var likeCount int

		err := rows.Scan(
			&blog.ID,
			&blog.UserId,
			&blog.Title,
			&blog.Description,
			&blog.GithubUrl,
			&blog.Category,
			&blog.Tags,
			&likeCount,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan blog: %v", err)
			return nil, err
		}

		// `Likes` フィールドにキャストして代入
		blog.Likes = int8(likeCount)

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

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
        SELECT b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags,
               COALESCE(COUNT(bl.id), 0) AS likes, b.created_at, b.updated_at
        FROM blogs b
        LEFT JOIN blogs_likes bl ON b.id = bl.blog_id
        WHERE b.id = $1
        GROUP BY b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags, b.created_at, b.updated_at
    `

	// Supabaseからクエリを実行し、条件に一致するデータを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, id)
	var likeCount int
	var blog models.BlogData
	err := row.Scan(
		&blog.ID,
		&blog.UserId,
		&blog.Title,
		&blog.Description,
		&blog.GithubUrl,
		&blog.Category,
		&blog.Tags,
		&likeCount,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to fetch blog: %v", err)
		return nil, err
	}

	// `Likes` フィールドにキャストして代入
	blog.Likes = int8(likeCount)

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

// ブログデータの更新
func (r *BlogRepositoryImpl) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	log.Printf("UpdateBlog start...")

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
        WITH updated_blog AS (
            UPDATE blogs
            SET title = $2, github_url = $3, category = $4, description = $5, tags = $6
            WHERE id = $1
            RETURNING id, user_id, title, description, github_url, category, tags, created_at, updated_at
        )
        SELECT ub.id, ub.user_id, ub.title, ub.description, ub.github_url, ub.category, ub.tags, 
               COALESCE(COUNT(bl.id), 0) AS likes, ub.created_at, ub.updated_at
        FROM updated_blog ub
        LEFT JOIN blogs_likes bl ON ub.id = bl.blog_id
        GROUP BY ub.id, ub.user_id, ub.title, ub.description, ub.github_url, ub.category, ub.tags, ub.created_at, ub.updated_at
    `

	// Supabaseからクエリを実行し、指定されたブログデータを更新
	row := supabase.Pool.QueryRow(supabase.Ctx, query, id, title, githubUrl, category, description, tags)

	// 結果をスキャンして更新されたブログデータを返す
	var likeCount int
	var blog models.BlogData
	err := row.Scan(
		&blog.ID,
		&blog.UserId,
		&blog.Title,
		&blog.Description,
		&blog.GithubUrl,
		&blog.Category,
		&blog.Tags,
		&likeCount,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	// `Likes` フィールドにキャストして代入
	blog.Likes = int8(likeCount)

	if err != nil {
		log.Printf("Failed to update blog: %v", err)
		return nil, err
	}

	log.Printf("Updated blog: %v", blog)
	return &blog, nil
}

// ブログデータの削除
func (r *BlogRepositoryImpl) DeleteBlog(id string) error {
	log.Printf("DeleteBlog start...")

	query := `
		DELETE FROM blogs
		WHERE id = $1
	`

	// Supabaseからクエリを実行し、指定されたブログデータを削除
	_, err := supabase.Pool.Exec(supabase.Ctx, query, id)
	if err != nil {
		log.Printf("Failed to delete blog: %v", err)
		return err
	}

	log.Println("Deleted blog successfully")
	return nil
}

// ブログカテゴリ一覧を取得する
func (r *BlogRepositoryImpl) FetchBlogCategories() ([]string, error) {
	log.Printf("FetchBlogCategories start...")

	query := `
		SELECT DISTINCT category
		FROM blogs
		ORDER BY category
	`

	// Supabaseからクエリを実行し、全カテゴリデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		log.Printf("Failed to fetch blog categories: %v", err)
		return nil, err
	}
	log.Println("Fetched blog categories successfully")
	defer rows.Close()

	var categories []string

	// 結果をスキャンしてカテゴリデータをリストに追加
	for rows.Next() {
		var category string

		err := rows.Scan(&category)
		if err != nil {
			log.Printf("Failed to scan category: %v", err)
			return nil, err
		}

		categories = append(categories, category)
	}

	if rows.Err() != nil {
		log.Printf("Failed to fetch blog categories: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Printf("Fetched %d blog categories", len(categories))
	return categories, nil
}
