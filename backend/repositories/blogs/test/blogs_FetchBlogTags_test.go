//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: タグ一覧を取得でき、シード済みの "go,testing" が含まれる。
func TestRepository_FetchBlogTags(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	tags, err := repo.FetchBlogTags()

	assert.NoError(t, err)
	assert.NotNil(t, tags)
	assert.NotEmpty(t, tags)
	assert.Contains(t, tags, "go,testing")
}
