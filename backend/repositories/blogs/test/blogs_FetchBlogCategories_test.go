//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: カテゴリ一覧を取得でき、シード済みの "Tech" が含まれる。
func TestRepository_FetchBlogCategories(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	categories, err := repo.FetchBlogCategories()

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.NotEmpty(t, categories)
	assert.Contains(t, categories, "Tech")
}
