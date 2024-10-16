package repositories_comments

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// ブログIDに一致するコメント情報を取得する
func (r *CommentRepositoryImpl) FetchCommentsByBlogId(blogId string) ([]models.CommentData, error) {
	log.Printf("FetchCommentsByBlogId start...")

	query := `
		SELECT id, blog_id, guest_user, comment, created_at
		FROM comments
		WHERE blog_id = $1
	`

	// Supabaseからクエリを実行し、条件に一致するデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query, blogId)
	if err != nil {
		log.Printf("Failed to fetch comments: %v", err)
		return nil, err
	}
	log.Println("Fetched comments successfully")
	defer rows.Close()

	var comments []models.CommentData

	// 結果をスキャンしてブログデータをリストに追加
	for rows.Next() {
		var comment models.CommentData
		err := rows.Scan(
			&comment.ID,
			&comment.BlogId,
			&comment.GuestUser,
			&comment.Comment,
			&comment.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan comment: %v", err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		log.Printf("Failed to fetch comments: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Printf("Fetched comments: %v", comments)
	return comments, nil
}

// コメント情報を新規作成する
func (r *CommentRepositoryImpl) CreateComment(blogId, guestUser, comment string) (*models.CommentData, error) {
	log.Printf("CreateComment start...")

	query := `
		INSERT INTO comments (blog_id, guest_user, comment)
		VALUES ($1, $2, $3)
		RETURNING id, blog_id, guest_user, comment, created_at
	`

	// Supabaseからクエリを実行し、新規作成したデータを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, blogId, guestUser, comment)
	var newComment models.CommentData
	err := row.Scan(
		&newComment.ID,
		&newComment.BlogId,
		&newComment.GuestUser,
		&newComment.Comment,
		&newComment.CreatedAt,
	)
	if err != nil {
		log.Printf("Failed to create comment: %v", err)
		return nil, err
	}

	log.Printf("Created comment: %v", newComment)
	return &newComment, nil
}
