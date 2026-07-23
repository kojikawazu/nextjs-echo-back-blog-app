//go:build integration

package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"os"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
)

// 正常系: 既知ユーザーのブログ一覧を取得でき、2件以上が返る。
func TestRepository_FetchBlogsByUserId(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	testUserId := os.Getenv("TEST_USER_ID")

	blogs, err := repo.FetchBlogsByUserId(testUserId)
	if err != nil {
		t.Fatalf("Failed to fetch blog: %v", err)
	}

	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.GreaterOrEqual(t, len(blogs), 2)
	for _, b := range blogs {
		assert.Equal(t, testUserId, b.BlogUserId)
	}
}

// 準正常系: UUID形式でないID（"2"）は 22P02 でエラーになる（SQLSTATE で検証）。
func TestRepository_FetchBlogsByUserId_ErrorCase(t *testing.T) {
	repo := repositories_blogs.NewBlogRepository()

	blogs, err := repo.FetchBlogsByUserId("2")

	assert.Error(t, err)
	assert.Nil(t, blogs)

	var pgErr *pgconn.PgError
	if assert.ErrorAs(t, err, &pgErr) {
		assert.Equal(t, "22P02", pgErr.Code)
	}
}
