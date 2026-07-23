//go:build integration

package repositories_blog_likes

import (
	"os"
	"testing"

	"backend/supabase"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// 異常系（安全な失敗）: コネクションプールがクローズ済みでも、
// いいね作成は panic せず error を返すこと。
func TestRepository_CreateBlogLike_ConnectionClosed(t *testing.T) {
	supabase.ClosePool()
	defer func() {
		if err := supabase.InitSupabase(); err != nil {
			t.Fatalf("failed to reinit pool: %v", err)
		}
	}()

	repo := NewBlogLikeRepository()

	like, err := repo.CreateBlogLike(os.Getenv("TEST_BLOG_ID"), uuid.New().String())

	assert.Error(t, err)
	assert.Nil(t, like)
}
