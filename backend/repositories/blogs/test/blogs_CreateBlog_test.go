//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 正常系: 作成したブログの created_at（timestamptz）・updated_at（timestamp）が現在時刻付近で返る。
// DB ドライバの移行（pgx v4 → v5, #153）で日時型の Scan が変わっていないことを確認する。
func TestRepository_CreateBlog_Timestamps(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	before := time.Now()
	blog, err := repo.CreateBlog(os.Getenv("TEST_USER_ID"), "test_title", "test_github_url", "test_category", "test_description", "test_tags")
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, repo.DeleteBlog(blog.ID))
	})

	assert.WithinDuration(t, before, blog.CreatedAt, time.Minute)
	assert.WithinDuration(t, before, blog.UpdatedAt, time.Minute)
}

func TestRepository_CreateBlog_Error(t *testing.T) {
	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// 異常系テスト
	blog, err := repo.CreateBlog("", "test_title", "test_github_url", "test_category", "test_description", "test_tags")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, blog)
}
