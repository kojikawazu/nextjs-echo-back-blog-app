package models

import (
	"fmt"
	"time"
)

// redacted はログ等の出力時に機密値の代わりに表示する文字列。
const redacted = "[REDACTED]"

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

// String は fmt の %v / %+v / %s 出力用の文字列を返す。
// ログへの機密情報の混入を構造的に防ぐため、Password と Email はマスクする。
//
// 戻り値:
//   - string: 機密値をマスクしたユーザー情報
func (u BlogUsersData) String() string {
	return fmt.Sprintf("BlogUsersData{ID:%s Name:%s Email:%s Password:%s CreatedAt:%s UpdatedAt:%s}",
		u.ID, u.Name, redacted, redacted, u.CreatedAt.Format(time.RFC3339), u.UpdatedAt.Format(time.RFC3339))
}

// GoString は fmt の %#v 出力用の文字列を返す。String と同様に機密値をマスクする。
//
// 戻り値:
//   - string: 機密値をマスクしたユーザー情報
func (u BlogUsersData) GoString() string {
	return "models." + u.String()
}
