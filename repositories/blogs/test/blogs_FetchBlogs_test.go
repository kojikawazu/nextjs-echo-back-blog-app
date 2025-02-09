package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchBlogs(t *testing.T) {
	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// メソッドを実行
	blogs, err := repo.FetchBlogs()
	if err != nil {
		t.Fatalf("Failed to fetch blogs: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.GreaterOrEqual(t, len(blogs), 2)
}
