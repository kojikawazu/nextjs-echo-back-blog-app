package repositories_comments

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchCommentsByBlogId(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewCommentRepository()

	// 環境変数から取得
	blogId := os.Getenv("TEST_BLOG_ID")

	// メソッドを実行
	comments, err := repo.FetchCommentsByBlogId(blogId)

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, comments)
}

func TestRepository_FetchCommentsByBlogId_InvalidBlogId(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewCommentRepository()

	// メソッドを実行
	comments, err := repo.FetchCommentsByBlogId("1")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, comments)
}
