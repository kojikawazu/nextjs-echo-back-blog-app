package models

import "time"

// BlogCommentsData はブログのコメント情報を表すデータ構造。
// API のレスポンスには dto.CommentResponse を使い、本構造体をそのまま返さない（JSON タグは DTO 側で管理する）。
type BlogCommentsData struct {
	ID        string    `db:"id"`         // UUID型
	BlogId    string    `db:"blog_id"`    // ブログID
	GuestUser string    `db:"guest_user"` // ゲスト名
	Comment   string    `db:"comment"`    // コメント
	CreatedAt time.Time `db:"created_at"` // タイムスタンプ
}
