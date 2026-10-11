//go:build integration

package repositories_blog_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// 正常系: 既知ユーザーを更新でき、更新後の値が返る。渡したハッシュがそのまま保存され、元のパスワードで照合できる。
func TestRepository_UpdateBlogUsers(t *testing.T) {
	repo := NewBlogUsersRepository()

	testId := os.Getenv("TEST_USER_ID")
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	// シードと同じパスワードの新しいハッシュで更新する（後続のテストのログインを壊さない）
	hash, err := bcrypt.GenerateFromPassword([]byte(testPasswd), bcrypt.DefaultCost)
	require.NoError(t, err)

	user, err := repo.UpdateBlogUsers(testId, testName, testEmail, string(hash))

	require.NoError(t, err)
	assert.Equal(t, testId, user.ID)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)

	stored, err := repo.FetchBlogUsersByEmail(testEmail)
	require.NoError(t, err)
	assert.Equal(t, string(hash), stored.Password)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(testPasswd)))
}

// 準正常系: 存在しないIDでの更新は該当行なし（pgx.ErrNoRows）を ErrUserNotFound に変換して返す。
func TestRepository_UpdateBlogUsers_NotFound(t *testing.T) {
	repo := NewBlogUsersRepository()

	user, err := repo.UpdateBlogUsers("00000000-0000-0000-0000-000000000000", "name", "email@example.com", "$2a$10$notarealhashnotarealhashnotarealhashnotarealhashnot")

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, user)
}
