//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
