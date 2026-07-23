//go:build integration

package repositories_blog_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: 既知ユーザーを更新でき、更新後の値が返る。
func TestRepository_UpdateBlogUsers(t *testing.T) {
	repo := NewBlogUsersRepository()

	testId := os.Getenv("TEST_USER_ID")
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	user, err := repo.UpdateBlogUsers(testId, testName, testEmail, testPasswd)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testId, user.ID)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

// 準正常系: 存在しないIDでの更新は該当行なし（no rows）でエラーになる。
func TestRepository_UpdateBlogUsers_NotFound(t *testing.T) {
	repo := NewBlogUsersRepository()

	user, err := repo.UpdateBlogUsers("00000000-0000-0000-0000-000000000000", "name", "email@example.com", "pass")

	assert.Error(t, err)
	assert.Nil(t, user)
}
