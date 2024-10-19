package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchBlogById(t *testing.T) {
	// Supabaseクライアントの初期化
	repositories_blogs.SetupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// 環境変数から取得
	id := os.Getenv("TEST_BLOG_ID")

	// メソッドを実行
	blog, err := repo.FetchBlogById(id)

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blog)
}

func TestRepository_FetchBlogById_ErrorCase(t *testing.T) {
	// Supabaseクライアントの初期化
	repositories_blogs.SetupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// メソッドを実行
	blog, err := repo.FetchBlogById("2")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "ERROR: invalid input syntax for type uuid: \"2\" (SQLSTATE 22P02)", err.Error())
}
