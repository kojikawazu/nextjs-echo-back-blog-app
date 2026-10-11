package models

import (
	"fmt"
	"time"
)

// redacted はログ等の出力時に機密値の代わりに表示する文字列。
const redacted = "[REDACTED]"

// BlogUsersData はブログユーザーの情報を表すデータ構造。
// API のレスポンスには dto.UserResponse を使い、本構造体をそのまま返さない。
// 誤ってシリアライズされてもパスワードが出力されないよう、Password は json:"-" とする。
type BlogUsersData struct {
	ID        string    `db:"id"`                // UUID型
	Name      string    `db:"name"`              // ユーザー名
	Email     string    `db:"email"`             // メールアドレス
	Password  string    `json:"-" db:"password"` // パスワード（シリアライズしない）
	CreatedAt time.Time `db:"created_at"`        // タイムスタンプ
	UpdatedAt time.Time `db:"updated_at"`        // タイムスタンプ
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
