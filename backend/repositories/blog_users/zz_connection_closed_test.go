//go:build integration

package repositories_blog_users

import (
	"os"
	"testing"

	"backend/supabase"

	"github.com/stretchr/testify/assert"
)

// 異常系（安全な失敗）: コネクションプールがクローズ済みの状態でも、
// クエリは panic せず error を返すこと。
//
// ファイル名を zz_ 始まりにしてパッケージ内で最後に実行されるようにし、
// 後続テストへの影響を避ける。念のため defer でプールを再初期化する。
func TestRepository_FetchUserById_ConnectionClosed(t *testing.T) {
	supabase.ClosePool()
	defer func() {
		if err := supabase.InitSupabase(); err != nil {
			t.Fatalf("failed to reinit pool: %v", err)
		}
	}()

	repo := NewBlogUsersRepository()

	user, err := repo.FetchBlogUsersById(os.Getenv("TEST_USER_ID"))

	assert.Error(t, err)
	assert.Nil(t, user)
}
