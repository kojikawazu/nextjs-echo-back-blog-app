package models

import "time"

// BlogUsersData はブログユーザーの情報を表すデータ構造。
// 各フィールドには、JSONおよびデータベースのタグを指定。
type BlogUsersData struct {
	ID        string    `json:"id" db:"id"`                 // UUID型
	Name      string    `json:"name" db:"name"`             // ユーザー名
	Email     string    `json:"email" db:"email"`           // メールアドレス
	Password  string    `json:"password" db:"password"`     // パスワード
	CreatedAt time.Time `json:"created_at" db:"created_at"` // タイムスタンプ
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // タイムスタンプ
}
