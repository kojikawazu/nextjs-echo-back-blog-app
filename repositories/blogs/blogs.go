package repositories_blogs

import (
	"backend/logger"
	"backend/models"
	"backend/supabase"
	"errors"

	"github.com/google/uuid"
)

// 全ブログデータを取得する
func (r *BlogRepositoryImpl) FetchBlogs() ([]models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogs start...")

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
		SELECT b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags,
				COALESCE(l.like_count, 0) AS likes,
				COALESCE(c.comment_count, 0) AS comment_cnt,
				b.created_at, b.updated_at
		FROM blogs b
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS like_count
			FROM blogs_likes
			GROUP BY blog_id
		) l ON b.id = l.blog_id
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY blog_id
		) c ON b.id = c.blog_id 
		ORDER BY b.created_at DESC
    `

	// Supabaseからクエリを実行し、全データ取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blogs: %v", err)
		return nil, err
	}
	logger.InfoLog.Println("Fetched blogs successfully")
	defer rows.Close()

	var blogs []models.BlogData

	// 結果をスキャンしてブログデータをリストに追加
	for rows.Next() {
		var blog models.BlogData
		var likeCount int
		var commentCnt int

		err := rows.Scan(
			&blog.ID,
			&blog.UserId,
			&blog.Title,
			&blog.Description,
			&blog.GithubUrl,
			&blog.Category,
			&blog.Tags,
			&likeCount,
			&commentCnt,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			logger.ErrorLog.Printf("Failed to scan blog: %v", err)
			return nil, err
		}

		// `Likes` フィールドにキャストして代入
		blog.Likes = int8(likeCount)
		blog.CommentCnt = int8(commentCnt)
		blogs = append(blogs, blog)
	}

	if rows.Err() != nil {
		logger.ErrorLog.Printf("Failed to fetch blogs: %v", rows.Err())
		return nil, rows.Err()
	}

	logger.InfoLog.Printf("Fetched %d blogs", len(blogs))
	return blogs, nil
}

// 指定されたユーザーIDに一致するブログデータを取得する
func (r *BlogRepositoryImpl) FetchBlogsByUserId(userId string) ([]models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogsByUserId start...")

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
		SELECT b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags, 
				COALESCE(l.like_count, 0) AS likes,
				COALESCE(c.comment_count, 0) AS comment_cnt,
				b.created_at, b.updated_at
		FROM blogs b
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS like_count
			FROM blogs_likes
			GROUP BY blog_id
		) l ON b.id = l.blog_id
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY blog_id
		) c ON b.id = c.blog_id 
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC
	`

	// Supabaseからクエリを実行し、条件に一致するデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query, userId)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blogs: %v", err)
		return nil, err
	}
	logger.InfoLog.Println("Fetched blogs successfully")
	defer rows.Close()

	var blogs []models.BlogData

	// 結果をスキャンしてブログデータをリストに追加
	for rows.Next() {
		var blog models.BlogData
		var likeCount int
		var commentCnt int

		err := rows.Scan(
			&blog.ID,
			&blog.UserId,
			&blog.Title,
			&blog.Description,
			&blog.GithubUrl,
			&blog.Category,
			&blog.Tags,
			&likeCount,
			&commentCnt,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			logger.ErrorLog.Printf("Failed to scan blog: %v", err)
			return nil, err
		}

		// `Likes` フィールドにキャストして代入
		blog.Likes = int8(likeCount)
		blog.CommentCnt = int8(commentCnt)

		blogs = append(blogs, blog)
	}

	if rows.Err() != nil {
		logger.ErrorLog.Printf("Failed to fetch blogs: %v", rows.Err())
		return nil, rows.Err()
	}

	logger.InfoLog.Printf("Fetched %d blogs", len(blogs))
	return blogs, nil
}

// 指定されたIDに一致するブログデータを取得する
func (r *BlogRepositoryImpl) FetchBlogById(id string) (*models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogById start...")

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
        SELECT b.id, b.user_id, b.title, b.description, b.github_url, b.category, b.tags,
				COALESCE(l.like_count, 0) AS likes,
				COALESCE(c.comment_count, 0) AS comment_cnt,
				b.created_at, b.updated_at
        FROM blogs b
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS like_count
			FROM blogs_likes
			GROUP BY blog_id
		) l ON b.id = l.blog_id
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY blog_id
		) c ON b.id = c.blog_id
        WHERE b.id = $1
    `

	// Supabaseからクエリを実行し、条件に一致するデータを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, id)
	var likeCount int
	var commentCnt int

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
		&commentCnt,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blog: %v", err)
		return nil, err
	}

	// `Likes` フィールドにキャストして代入
	blog.Likes = int8(likeCount)
	blog.CommentCnt = int8(commentCnt)

	logger.InfoLog.Printf("Fetched blog: %v", blog)
	return &blog, nil
}

// ブログデータの作成
func (r *BlogRepositoryImpl) CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	logger.InfoLog.Printf("CreateBlog start...")

	if userId == "" {
		return nil, errors.New("user_id cannot be empty")
	}
	if _, err := uuid.Parse(userId); err != nil {
		return nil, errors.New("invalid user_id format: must be a valid UUID")
	}

	query := `
		INSERT INTO blogs (user_id, title, github_url, category, description, tags)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, description, github_url, category, tags, likes, comment_cnt, created_at, updated_at
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
		&blog.CommentCnt,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)
	if err != nil {
		logger.ErrorLog.Printf("Failed to create blog: %v", err)
		return nil, err
	}

	logger.InfoLog.Printf("Created blog: %v", blog)
	return &blog, nil
}

// ブログデータの更新
func (r *BlogRepositoryImpl) UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error) {
	logger.InfoLog.Printf("UpdateBlog start...")

	// ※ blogs と blogs_likes テーブルを結合し、いいね数を集計して取得すること
	query := `
        WITH updated_blog AS (
            UPDATE blogs
            SET title = $2, github_url = $3, category = $4, description = $5, tags = $6
            WHERE id = $1
            RETURNING id, user_id, title, description, github_url, category, tags, created_at, updated_at
        )
        SELECT ub.id, ub.user_id, ub.title, ub.description, ub.github_url, ub.category, ub.tags, 
               COALESCE(l.like_count, 0) AS likes,
			   COALESCE(c.comment_count, 0) AS comment_cnt,
			   ub.created_at, ub.updated_at
        FROM updated_blog ub
        LEFT JOIN (
			SELECT blog_id, COUNT(*) AS like_count
			FROM blogs_likes
			GROUP BY blog_id
		) l ON ub.id = l.blog_id
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS comment_count
			FROM comments
			GROUP BY blog_id
		) c ON ub.id = c.blog_id
    `

	// Supabaseからクエリを実行し、指定されたブログデータを更新
	row := supabase.Pool.QueryRow(supabase.Ctx, query, id, title, githubUrl, category, description, tags)

	// 結果をスキャンして更新されたブログデータを返す
	var likeCount int
	var commentCnt int
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
		&commentCnt,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	// キャストして代入
	blog.Likes = int8(likeCount)
	blog.CommentCnt = int8(commentCnt)
	if err != nil {
		logger.ErrorLog.Printf("Failed to update blog: %v", err)
		return nil, err
	}

	logger.InfoLog.Printf("Updated blog: %v", blog)
	return &blog, nil
}

// ブログデータの削除
func (r *BlogRepositoryImpl) DeleteBlog(id string) error {
	logger.InfoLog.Printf("DeleteBlog start...")

	if id == "" {
		return errors.New("id cannot be empty")
	}
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid id format: must be a valid UUID")
	}

	query := `
		DELETE FROM blogs
		WHERE id = $1
	`

	// Supabaseからクエリを実行し、指定されたブログデータを削除
	_, err := supabase.Pool.Exec(supabase.Ctx, query, id)
	if err != nil {
		logger.ErrorLog.Printf("Failed to delete blog: %v", err)
		return err
	}

	logger.InfoLog.Println("Deleted blog successfully")
	return nil
}

// ブログカテゴリ一覧を取得する
func (r *BlogRepositoryImpl) FetchBlogCategories() ([]string, error) {
	logger.InfoLog.Printf("FetchBlogCategories start...")

	query := `
		SELECT DISTINCT category
		FROM blogs
		ORDER BY category
	`

	// Supabaseからクエリを実行し、全カテゴリデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blog categories: %v", err)
		return nil, err
	}
	logger.InfoLog.Println("Fetched blog categories successfully")
	defer rows.Close()

	var categories []string

	// 結果をスキャンしてカテゴリデータをリストに追加
	for rows.Next() {
		var category string

		err := rows.Scan(&category)
		if err != nil {
			logger.ErrorLog.Printf("Failed to scan category: %v", err)
			return nil, err
		}

		categories = append(categories, category)
	}

	if rows.Err() != nil {
		logger.ErrorLog.Printf("Failed to fetch blog categories: %v", rows.Err())
		return nil, rows.Err()
	}

	logger.InfoLog.Printf("Fetched %d blog categories", len(categories))
	return categories, nil
}

// ブログタグ一覧を取得する
func (r *BlogRepositoryImpl) FetchBlogTags() ([]string, error) {
	logger.InfoLog.Printf("FetchBlogTags start...")

	query := `
		SELECT DISTINCT tags
		FROM blogs
		ORDER BY tags
	`

	// Supabaseからクエリを実行し、全タグデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch blog tags: %v", err)
		return nil, err
	}
	logger.InfoLog.Println("Fetched blog tags successfully")
	defer rows.Close()

	var tags []string

	// 結果をスキャンしてタグデータをリストに追加
	for rows.Next() {
		var tag string

		err := rows.Scan(&tag)
		if err != nil {
			logger.ErrorLog.Printf("Failed to scan tag: %v", err)
			return nil, err
		}

		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		logger.ErrorLog.Printf("Failed to fetch blog tags: %v", rows.Err())
		return nil, rows.Err()
	}

	logger.InfoLog.Printf("Fetched %d blog tags", len(tags))
	return tags, nil
}

// 人気のあるブログを取得する
func (r *BlogRepositoryImpl) FetchBlogPopular(count int) ([]models.BlogData, error) {
	logger.InfoLog.Printf("FetchBlogPopular start...")

	query := `
		SELECT b.id, b.user_id, b.title, 
			   COALESCE(l.like_count, 0) AS likes,
			   b.created_at, b.updated_at
		FROM blogs b
		LEFT JOIN (
			SELECT blog_id, COUNT(*) AS like_count
			FROM blogs_likes
			GROUP BY blog_id
		) l ON b.id = l.blog_id
		ORDER BY likes DESC
		LIMIT $1
	`

	// Supabaseからクエリを実行し、人気のあるブログデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query, count)
	if err != nil {
		logger.ErrorLog.Printf("Failed to fetch popular blogs: %v", err)
		return nil, err
	}
	logger.InfoLog.Println("Fetched popular blogs successfully")
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
			&likeCount,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			logger.ErrorLog.Printf("Failed to scan blog: %v", err)
			return nil, err
		}

		// `Likes` フィールドにキャストして代入
		blog.Likes = int8(likeCount)

		// 空のものも設定しておく
		blog.Description = ""
		blog.GithubUrl = ""
		blog.Category = ""
		blog.Tags = ""

		blogs = append(blogs, blog)
	}

	if rows.Err() != nil {
		logger.ErrorLog.Printf("Failed to fetch popular blogs: %v", rows.Err())
		return nil, rows.Err()
	}

	logger.InfoLog.Printf("Fetched %d popular blogs", len(blogs))
	return blogs, nil
}
