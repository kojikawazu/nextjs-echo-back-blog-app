//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"

	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: 全ブログを取得でき、シード済みの2件以上が返る。
func TestRepository_FetchBlogs(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	blogs, err := repo.FetchBlogs()
	if err != nil {
		t.Fatalf("Failed to fetch blogs: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.GreaterOrEqual(t, len(blogs), 2)

	// シード済みのタイトルが含まれることを具体的に検証する。
	titles := make([]string, 0, len(blogs))
	for _, b := range blogs {
		titles = append(titles, b.Title)
	}
	assert.Contains(t, titles, "Test Blog Title")
}
