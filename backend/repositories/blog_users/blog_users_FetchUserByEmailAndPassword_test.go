//go:build integration

package repositories_blog_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: 既知のメール・パスワードでユーザーを取得できる。
func TestRepository_FetchUserByEmailAndPassword(t *testing.T) {
	repo := NewBlogUsersRepository()

	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	user, err := repo.FetchBlogUsersByEmailAndPassword(testEmail, testPasswd)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

// 準正常系: 空のメール・パスワードは該当行なしでエラーになり、ユーザーは返らない。
func TestRepository_FetchUserByEmailAndPassword_ErrorCases(t *testing.T) {
	repo := NewBlogUsersRepository()

	user, err := repo.FetchBlogUsersByEmailAndPassword("", "")

	assert.Error(t, err)
	assert.Nil(t, user)
}
