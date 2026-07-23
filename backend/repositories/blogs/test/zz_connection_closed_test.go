//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"backend/supabase"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 異常系（安全な失敗）: コネクションプールがクローズ済みでも、
// ブログ一覧取得は panic せず error を返すこと。
//
// ファイル名を zz_ 始まりにしてパッケージ内で最後に実行されるようにし、
// 後続テストへの影響を避ける。念のため defer でプールを再初期化する。
func TestRepository_FetchBlogs_ConnectionClosed(t *testing.T) {
	supabase.ClosePool()
	defer func() {
		if err := supabase.InitSupabase(); err != nil {
			t.Fatalf("failed to reinit pool: %v", err)
		}
	}()

	repo := repositories_blogs.NewBlogRepository()

	blogs, err := repo.FetchBlogs()

	assert.Error(t, err)
	assert.Nil(t, blogs)
}
