package repositories_blog_comments

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// FetchCommentsByBlogId はブログIDに一致するコメント情報を取得する。
//
// 引数:
//   - blogId: 取得対象のブログID
//
// 戻り値:
//   - []models.BlogCommentsData: 取得したコメント情報の一覧
//   - error: 取得に失敗した場合のエラー
func (r *CommentRepositoryImpl) FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error) {
	log.Printf("FetchCommentsByBlogId start...")

	query := `
		SELECT id, blog_id, guest_user, comment, created_at
		FROM blog_comments
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

	var comments []models.BlogCommentsData

	// 結果をスキャンしてブログデータをリストに追加
	for rows.Next() {
		var comment models.BlogCommentsData
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

// CreateComment はコメント情報を新規作成する。
//
// 引数:
//   - blogId: コメント対象のブログID
//   - guestUser: コメントを投稿したゲストユーザー名
//   - comment: コメント本文
//
// 戻り値:
//   - *models.BlogCommentsData: 作成したコメント情報
//   - error: 作成に失敗した場合のエラー
func (r *CommentRepositoryImpl) CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error) {
	log.Printf("CreateComment start...")

	query := `
		INSERT INTO blog_comments (blog_id, guest_user, comment)
		VALUES ($1, $2, $3)
		RETURNING id, blog_id, guest_user, comment, created_at
	`

	// Supabaseからクエリを実行し、新規作成したデータを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, blogId, guestUser, comment)
	var newComment models.BlogCommentsData
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
