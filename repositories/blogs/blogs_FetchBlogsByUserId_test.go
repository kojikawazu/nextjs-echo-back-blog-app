package repositories_blogs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchBlogsByUserId(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewBlogRepository()

	// 環境変数からユーザIDを取得
	testUserId := os.Getenv("TEST_USER_ID")

	// メソッドを実行
	blogs, err := repo.FetchBlogsByUserId(testUserId)
	if err != nil {
		t.Fatalf("Failed to fetch blog: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.Len(t, blogs, 1)
}

func TestRepository_FetchBlogsByUserId_ErrorCase(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewBlogRepository()

	// メソッドを実行
	blogs, err := repo.FetchBlogsByUserId("2")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, blogs)
	assert.Equal(t, "ERROR: invalid input syntax for type uuid: \"2\" (SQLSTATE 22P02)", err.Error())
}
