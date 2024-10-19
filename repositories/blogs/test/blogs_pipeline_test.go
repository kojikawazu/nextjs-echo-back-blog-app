package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Blog_PipeLine(t *testing.T) {
	// Supabaseクライアントの初期化
	repositories_blogs.SetupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// 環境変数から取得
	userId := os.Getenv("TEST_USER_ID")

	// ----------------------------------------------------------------------------------------------------------------------------
	// 1. ブログ生成テスト
	// ----------------------------------------------------------------------------------------------------------------------------
	blog, err := repo.CreateBlog(userId, "test_title", "test_github_url", "test_category", "test_description", "test_tags")

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blog)

	// ----------------------------------------------------------------------------------------------------------------------------
	// 2. ブログ取得テスト
	// ----------------------------------------------------------------------------------------------------------------------------
	fetchedBlog, err := repo.FetchBlogById(blog.ID)

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, fetchedBlog)

	// ----------------------------------------------------------------------------------------------------------------------------
	// 3. ブログ更新テスト
	// ----------------------------------------------------------------------------------------------------------------------------
	updatedBlog, err := repo.UpdateBlog(blog.ID, "updated_title", "updated_github_url", "updated_category", "updated_description", "updated_tags")

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, updatedBlog)

	// ----------------------------------------------------------------------------------------------------------------------------
	// 4. ブログ削除テスト
	// ----------------------------------------------------------------------------------------------------------------------------
	err = repo.DeleteBlog(blog.ID)

	// エラーチェック
	assert.NoError(t, err)

	// ----------------------------------------------------------------------------------------------------------------------------
	// 5. ブログ取得エラーテスト
	// ----------------------------------------------------------------------------------------------------------------------------
	fetchedBlog, err = repo.FetchBlogById(blog.ID)

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, fetchedBlog)
}
