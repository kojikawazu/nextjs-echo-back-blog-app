package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchBlogPopular(t *testing.T) {
	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// メソッドを実行
	blogs, err := repo.FetchBlogPopular(10)

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.NotEmpty(t, blogs)
}

func TestRepository_FetchBlogPopular_Empty(t *testing.T) {
	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// メソッドを実行
	blogs, err := repo.FetchBlogPopular(0)

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.Nil(t, blogs)
	assert.Empty(t, blogs)
}
