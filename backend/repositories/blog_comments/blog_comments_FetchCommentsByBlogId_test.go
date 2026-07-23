//go:build integration

package repositories_blog_comments

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: 既知ブログのコメント一覧を取得でき、シード済みコメントが含まれる。
func TestRepository_FetchCommentsByBlogId(t *testing.T) {
	repo := NewCommentRepository()

	blogId := os.Getenv("TEST_BLOG_ID")

	comments, err := repo.FetchCommentsByBlogId(blogId)

	assert.NoError(t, err)
	assert.NotNil(t, comments)
	assert.GreaterOrEqual(t, len(comments), 1)
	for _, c := range comments {
		assert.Equal(t, blogId, c.BlogId)
	}
}

// 準正常系: UUID形式でないblogId（"1"）はエラーになり、一覧は返らない。
func TestRepository_FetchCommentsByBlogId_InvalidBlogId(t *testing.T) {
	repo := NewCommentRepository()

	comments, err := repo.FetchCommentsByBlogId("1")

	assert.Error(t, err)
	assert.Nil(t, comments)
}
