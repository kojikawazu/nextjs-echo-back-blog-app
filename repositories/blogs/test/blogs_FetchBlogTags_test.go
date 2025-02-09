package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchBlogTags(t *testing.T) {
	// リポジトリのインスタンスを作成
	repo := repositories_blogs.NewBlogRepository()

	// メソッドを実行
	tags, err := repo.FetchBlogTags()

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.NotEmpty(t, tags)
}
