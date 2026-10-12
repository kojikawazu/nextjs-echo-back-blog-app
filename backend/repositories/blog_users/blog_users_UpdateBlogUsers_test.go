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

// 正常系: 更新すると updated_at が DB のトリガで進み、created_at は変わらない（#117）。
// リポジトリの UPDATE 文は updated_at を書かないため、値が進むのはトリガによる。
func TestRepository_UpdateBlogUsers_UpdatedAt(t *testing.T) {
	repo := NewBlogUsersRepository()

	before, err := repo.FetchBlogUsersById(os.Getenv("TEST_USER_ID"))
	require.NoError(t, err)

	// 名前・メール・パスワードは元の値のまま更新する（後続のテストのログインを壊さない）
	updated, err := repo.UpdateBlogUsers(before.ID, before.Name, before.Email, before.Password)
	require.NoError(t, err)

	assert.True(t, updated.UpdatedAt.After(before.UpdatedAt),
		"updated_at が進んでいない: before=%v updated=%v", before.UpdatedAt, updated.UpdatedAt)
	assert.True(t, updated.CreatedAt.Equal(before.CreatedAt),
		"created_at が変わった: before=%v updated=%v", before.CreatedAt, updated.CreatedAt)
}

// 準正常系: 存在しないIDでの更新は該当行なし（pgx.ErrNoRows）を ErrUserNotFound に変換して返す。
func TestRepository_UpdateBlogUsers_NotFound(t *testing.T) {
	repo := NewBlogUsersRepository()

	user, err := repo.UpdateBlogUsers("00000000-0000-0000-0000-000000000000", "name", "email@example.com", "$2a$10$notarealhashnotarealhashnotarealhashnotarealhashnot")

	assert.ErrorIs(t, err, ErrUserNotFound)
	assert.Nil(t, user)
}
