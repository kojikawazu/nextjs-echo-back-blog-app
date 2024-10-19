package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_DeleteBlog_Error(t *testing.T) {
	// Supabaseクライアントの初期化
	repositories_blogs.SetupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// 異常系テスト
	err := repo.DeleteBlog("")

	// エラーチェックとデータ確認
	assert.Error(t, err)
}
