package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateBlog_Error(t *testing.T) {
	// Supabaseクライアントの初期化
	repositories_blogs.SetupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// 異常系テスト
	blog, err := repo.CreateBlog("", "test_title", "test_github_url", "test_category", "test_description", "test_tags")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Empty(t, blog.ID)
}
