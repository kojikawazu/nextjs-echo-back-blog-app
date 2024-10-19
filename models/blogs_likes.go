package models

import "time"

// ブログのいいね情報を表すデータ構造
// 各フィールドには、JSONおよびデータベースのタグを指定。
type BlogLikeData struct {
	ID        string    `json:"id" db:"id"`                 // UUID型
	BlogId    string    `json:"blog_id" db:"blog_id"`       // ブログID
	VisitId   string    `json:"visit_id" db:"visit_id"`     // 訪問ID
	CreatedAt time.Time `json:"created_at" db:"created_at"` // タイムスタンプ
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // タイムスタンプ
}
