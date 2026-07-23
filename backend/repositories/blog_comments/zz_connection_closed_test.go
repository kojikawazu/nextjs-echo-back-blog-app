//go:build integration

package repositories_blog_comments

import (
	"os"
	"testing"

	"backend/supabase"

	"github.com/stretchr/testify/assert"
)

// 異常系（安全な失敗）: コネクションプールがクローズ済みでも、
// コメント一覧取得は panic せず error を返すこと。
func TestRepository_FetchCommentsByBlogId_ConnectionClosed(t *testing.T) {
	supabase.ClosePool()
	defer func() {
		if err := supabase.InitSupabase(); err != nil {
			t.Fatalf("failed to reinit pool: %v", err)
		}
	}()

	repo := NewCommentRepository()

	comments, err := repo.FetchCommentsByBlogId(os.Getenv("TEST_BLOG_ID"))

	assert.Error(t, err)
	assert.Nil(t, comments)
}
