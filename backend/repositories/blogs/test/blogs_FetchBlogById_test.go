//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"os"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
)

// 正常系: 既知のブログIDでブログを取得でき、シード済みの値が返る。
func TestRepository_FetchBlogById(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	id := os.Getenv("TEST_BLOG_ID")

	blog, err := repo.FetchBlogById(id)

	assert.NoError(t, err)
	assert.NotNil(t, blog)
	assert.Equal(t, id, blog.ID)
	assert.Equal(t, os.Getenv("TEST_USER_ID"), blog.BlogUserId)
	assert.Equal(t, "Test Blog Title", blog.Title)
}

// 準正常系: UUID形式でないID（"2"）はPostgreSQLが 22P02 でエラーを返す。
// エラーメッセージ全文ではなく SQLSTATE コードで検証し、PG/ドライバ差異に強くする。
func TestRepository_FetchBlogById_ErrorCase(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	blog, err := repo.FetchBlogById("2")

	assert.Error(t, err)
	assert.Nil(t, blog)

	var pgErr *pgconn.PgError
	if assert.ErrorAs(t, err, &pgErr) {
		assert.Equal(t, "22P02", pgErr.Code)
	}
}
