package models

import (
	"time"
)

// BlogData はブログの情報を表すデータ構造。
// API のレスポンスには dto.BlogResponse を使い、本構造体をそのまま返さない（JSON タグは DTO 側で管理する）。
type BlogData struct {
	ID          string    `db:"id"`           // UUID型
	BlogUserId  string    `db:"blog_user_id"` // ブログユーザーID
	Title       string    `db:"title"`        // タイトル
	Description string    `db:"description"`  // 説明
	GithubUrl   string    `db:"github_url"`   // GitHubリポジトリのURL
	Category    string    `db:"category"`     // カテゴリ
	Tags        string    `db:"tags"`         // タグ
	Likes       int8      `db:"likes"`        // いいね数
	CommentCnt  int8      `db:"comment_cnt"`  // コメント数
	CreatedAt   time.Time `db:"created_at"`   // タイムスタンプ
	UpdatedAt   time.Time `db:"updated_at"`   // タイムスタンプ
}
