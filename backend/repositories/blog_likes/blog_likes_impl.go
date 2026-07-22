package repositories_blog_likes

import (
	"backend/models"
	"backend/supabase"
	"log"
)

// FetchBlogLikesByVisitId は VisitId によっていいねデータを取得する。
func (r *BlogLikeRepositoryImpl) FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error) {
	log.Println("FetchBlogLikesByVisitId start...")

	// データベースからいいねデータを取得
	query := `
		SELECT id, blog_id, visit_id, created_at, updated_at
		FROM blog_likes
		WHERE visit_id = $1
	`
	// クエリを実行し、いいねデータを取得
	rows, err := supabase.Pool.Query(supabase.Ctx, query, visitId)
	if err != nil {
		log.Printf("Failed to fetch blog likes by visit id: %v", err)
		return nil, err
	}
	defer rows.Close()

	// いいねデータを格納するスライスを作成
	var blogLikes []models.BlogLikesData

	// いいねデータをスキャンしてスライスに追加
	for rows.Next() {
		var blogLike = models.BlogLikesData{}
		err := rows.Scan(&blogLike.ID, &blogLike.BlogId, &blogLike.VisitId, &blogLike.CreatedAt, &blogLike.UpdatedAt)
		if err != nil {
			log.Printf("Failed to scan blog like: %v", err)
			return nil, err
		}
		blogLikes = append(blogLikes, blogLike)
	}

	log.Printf("Fetched blog likes by visit id: %v", blogLikes)
	return blogLikes, nil
}

// IsBlogLiked はいいねが存在するか確認する。
func (r *BlogLikeRepositoryImpl) IsBlogLiked(blogId, visitId string) (bool, error) {
	log.Println("IsBlogLiked start...")

	// データベースからいいねデータを取得
	query := `
		SELECT id
		FROM blog_likes
		WHERE blog_id = $1 AND visit_id = $2
	`
	// クエリを実行し、いいねデータを取得
	row := supabase.Pool.QueryRow(supabase.Ctx, query, blogId, visitId)
	var id string

	// スキャンしていいねデータが存在するか確認
	err := row.Scan(&id)
	if err != nil {
		log.Printf("Failed to check if blog is liked: %v", err)
		return false, err
	}

	log.Println("Blog is liked")
	return true, nil
}

// CreateBlogLike はいいねデータを作成する。
func (r *BlogLikeRepositoryImpl) CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error) {
	log.Println("CreateBlogLike start...")

	// データベースにいいねデータを挿入
	query := `
		INSERT INTO blog_likes (blog_id, visit_id)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	// クエリを実行し、新しいいいねデータを作成
	row := supabase.Pool.QueryRow(supabase.Ctx, query, blogId, visitId)
	var blogLike = &models.BlogLikesData{}

	// スキャンしていいねデータを返す
	err := row.Scan(&blogLike.ID, &blogLike.CreatedAt, &blogLike.UpdatedAt)
	if err != nil {
		log.Printf("Failed to create blog like: %v", err)
		return nil, err
	}

	log.Printf("Created blog like: %v", blogLike)
	return blogLike, nil
}

// DeleteBlogLike はいいねデータを削除する。
func (r *BlogLikeRepositoryImpl) DeleteBlogLike(blogId, visitId string) error {
	log.Println("DeleteBlogLike start...")

	// データベースからいいねデータを削除
	query := `
		DELETE FROM blog_likes
		WHERE blog_id = $1 AND visit_id = $2
	`
	// クエリを実行し、いいねデータを削除
	_, err := supabase.Pool.Exec(supabase.Ctx, query, blogId, visitId)
	if err != nil {
		log.Printf("Failed to delete blog like: %v", err)
		return err
	}

	log.Println("Deleted blog like")
	return nil
}
