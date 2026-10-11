//go:build integration

package repositories_blog_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// 正常系: 既知のメールアドレスでユーザーを取得でき、パスワードは平文ではなく bcrypt ハッシュで返る。
func TestRepository_FetchUserByEmail(t *testing.T) {
	repo := NewBlogUsersRepository()

	testID := os.Getenv("TEST_USER_ID")
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	user, err := repo.FetchBlogUsersByEmail(testEmail)
	require.NoError(t, err)

	assert.Equal(t, testID, user.ID)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
	assert.NotEqual(t, testPasswd, user.Password, "DB 上のパスワードは平文と一致しない")
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testPasswd)))
}

// 準正常系: 存在しないメールアドレスは ErrUserNotFound を返す（#137: pgx.ErrNoRows を変換する）。
func TestRepository_FetchUserByEmail_NotFound(t *testing.T) {
	repo := NewBlogUsersRepository()

	for _, email := range []string{"notexist@example.com", ""} {
		user, err := repo.FetchBlogUsersByEmail(email)

		assert.ErrorIs(t, err, ErrUserNotFound, "email=%q", email)
		assert.Nil(t, user)
	}
}
