package models

import "time"

// BlogLikesData はブログのいいね情報を表すデータ構造。
// API のレスポンスには dto.BlogLikeResponse を使い、本構造体をそのまま返さない（JSON タグは DTO 側で管理する）。
type BlogLikesData struct {
	ID        string    `db:"id"`         // UUID型
	BlogId    string    `db:"blog_id"`    // ブログID
	VisitId   string    `db:"visit_id"`   // 訪問ID
	CreatedAt time.Time `db:"created_at"` // タイムスタンプ
	UpdatedAt time.Time `db:"updated_at"` // タイムスタンプ
}
