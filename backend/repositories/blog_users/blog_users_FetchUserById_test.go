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

// 準正常系: 空文字IDはUUID変換不可でエラーになり、ユーザーは返らない。
func TestRepository_FetchUserById_InvalidId(t *testing.T) {
	repo := NewBlogUsersRepository()

	user, err := repo.FetchBlogUsersById("")

	assert.Error(t, err)
	assert.Nil(t, user)
}
