package models

import "time"

// ブログのコメント情報を表すデータ構造
// 各フィールドには、JSONおよびデータベースのタグを指定。
type CommentData struct {
	ID        string    `json:"id" db:"id"`                 // UUID型
	BlogId    string    `json:"blog_id" db:"blog_id"`       // ブログID
	GuestUser string    `json:"guest_user" db:"guest_user"` // ゲスト名
	Comment   string    `json:"comment" db:"comment"`       // コメント
	CreatedAt time.Time `json:"created_at" db:"created_at"` // タイムスタンプ
}
