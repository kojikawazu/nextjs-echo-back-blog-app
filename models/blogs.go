package models

import (
	"time"
)

// ブログの情報を表すデータ構造
// 各フィールドには、JSONおよびデータベースのタグを指定。
type BlogData struct {
	ID          string    `json:"id" db:"id"`                   // UUID型
	UserId      string    `json:"user_id" db:"user_id"`         // ユーザーID
	Title       string    `json:"title" db:"title"`             // タイトル
	Description string    `json:"description" db:"description"` // 説明
	GithubUrl   string    `json:"github_url" db:"github_url"`   // GitHubリポジトリのURL
	Category    string    `json:"category" db:"category"`       // カテゴリ
	Tag         string    `json:"tag" db:"tag"`                 // タグ
	Likes       int8      `json:"likes" db:"likes"`             // いいね数
	CreatedAt   time.Time `json:"created_at" db:"created_at"`   // タイムスタンプ
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`   // タイムスタンプ
}
