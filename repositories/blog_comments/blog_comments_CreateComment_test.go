package repositories_blog_comments

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateComment(t *testing.T) {
	setupSupabase(t)

	repo := NewCommentRepository()

	blogId := os.Getenv("TEST_BLOG_ID")

	comment, err := repo.CreateComment(blogId, "guest_user", "test comment")

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, blogId, comment.BlogId)
	assert.Equal(t, "guest_user", comment.GuestUser)
	assert.Equal(t, "test comment", comment.Comment)
	assert.NotEmpty(t, comment.ID)
}

func TestRepository_CreateComment_InvalidBlogId(t *testing.T) {
	setupSupabase(t)

	repo := NewCommentRepository()

	// 不正なblogId（UUID形式でない）→ PostgreSQLがエラー
	comment, err := repo.CreateComment("invalid-blog-id", "guest", "comment")

	assert.Error(t, err)
	assert.Nil(t, comment)
}
