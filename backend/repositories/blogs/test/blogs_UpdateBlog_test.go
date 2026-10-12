//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 正常系: 更新すると updated_at が DB のトリガで進み、created_at は変わらない（#117）。
// リポジトリの UPDATE 文は updated_at を書かないため、値が進むのはトリガによる。
func TestRepository_UpdateBlog_UpdatedAt(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	created, err := repo.CreateBlog(os.Getenv("TEST_USER_ID"), "test_title", "test_github_url", "test_category", "test_description", "test_tags")
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, repo.DeleteBlog(created.ID))
	})

	updated, err := repo.UpdateBlog(created.ID, "updated_title", "updated_github_url", "updated_category", "updated_description", "updated_tags")
	require.NoError(t, err)

	assert.True(t, updated.UpdatedAt.After(created.UpdatedAt),
		"updated_at が進んでいない: created=%v updated=%v", created.UpdatedAt, updated.UpdatedAt)
	assert.True(t, updated.CreatedAt.Equal(created.CreatedAt),
		"created_at が変わった: created=%v updated=%v", created.CreatedAt, updated.CreatedAt)
}

func TestRepository_UpdateBlog_Error(t *testing.T) {
	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// 異常系テスト
	updatedBlog, err := repo.UpdateBlog("", "updated_title", "updated_github_url", "updated_category", "updated_description", "updated_tags")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, updatedBlog)
}

// 準正常系: 存在しない UUID の場合は pgx.ErrNoRows ではなく ErrBlogNotFound を返す（ドライバ固有のエラーを漏らさない）。
func TestRepository_UpdateBlog_NotFound(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	updatedBlog, err := repo.UpdateBlog(notExistBlogID, "updated_title", "updated_github_url", "updated_category", "updated_description", "updated_tags")

	assert.ErrorIs(t, err, repositories_blogs.ErrBlogNotFound)
	assert.Nil(t, updatedBlog)
}
