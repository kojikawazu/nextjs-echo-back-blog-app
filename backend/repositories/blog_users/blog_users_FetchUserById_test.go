//go:build integration

package repositories_blog_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: 既知のユーザーIDでユーザーを取得できる。
func TestRepository_FetchUserById(t *testing.T) {
	repo := NewBlogUsersRepository()

	testUserId := os.Getenv("TEST_USER_ID")
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")

	user, err := repo.FetchBlogUsersById(testUserId)
	if err != nil {
		t.Fatalf("Failed to fetch user by id: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, testUserId, user.ID)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

// 準正常系: 存在しないIDは該当行なし（pgx.ErrNoRows）を ErrUserNotFound に変換して返す。
func TestRepository_FetchUserById_NotFound(t *testing.T) {
	repo := NewBlogUsersRepository()

	user, err := repo.FetchBlogUsersById("00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, user)
}

// 異常系: 空文字IDはUUID変換不可でエラーになり、ユーザーは返らない。不在（ErrUserNotFound）とは区別する。
func TestRepository_FetchUserById_InvalidId(t *testing.T) {
	repo := NewBlogUsersRepository()

	user, err := repo.FetchBlogUsersById("")

	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, user)
}
