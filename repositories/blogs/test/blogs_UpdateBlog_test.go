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
