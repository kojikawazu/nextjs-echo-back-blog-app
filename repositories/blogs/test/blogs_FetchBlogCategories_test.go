package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchBlogCategories(t *testing.T) {
	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// メソッドを実行
	categories, err := repo.FetchBlogCategories()

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.NotEmpty(t, categories)
}
