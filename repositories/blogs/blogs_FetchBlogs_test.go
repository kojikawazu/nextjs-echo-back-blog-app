package repositories_blogs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchBlogs(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewBlogRepository()

	// メソッドを実行
	blogs, err := repo.FetchBlogs()
	if err != nil {
		t.Fatalf("Failed to fetch blogs: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.Len(t, blogs, 2)
}
