package services_blog_users

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mustHash はテスト用にパスワードの bcrypt ハッシュを作る。
func mustHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := hashPassword(password)
	require.NoError(t, err)
	return hash
}

// isHashOf は渡された値が password の bcrypt ハッシュであるかを判定する mock.MatchedBy 用の述語を返す。
func isHashOf(password string) func(string) bool {
	return func(hash string) bool {
		return hash != password && verifyPassword(hash, password)
	}
}

// 正常系: ハッシュは平文と異なり、元のパスワードで照合でき、毎回異なるソルトで生成される。
func TestHashPassword(t *testing.T) {
	hash, err := hashPassword("correct-password")

	require.NoError(t, err)
	assert.NotEqual(t, "correct-password", hash)
	assert.True(t, strings.HasPrefix(hash, "$2a$10$"), "bcrypt cost 10 のハッシュであること: %s", hash)
	assert.True(t, verifyPassword(hash, "correct-password"))

	hash2, err := hashPassword("correct-password")
	require.NoError(t, err)
	assert.NotEqual(t, hash, hash2, "ソルトにより同じパスワードでも別のハッシュになる")
}

// 準正常系: 72 バイトちょうどはハッシュ化でき、73 バイトは errPasswordTooLong になる。
func TestHashPassword_Length(t *testing.T) {
	_, err := hashPassword(strings.Repeat("a", 72))
	assert.NoError(t, err)

	hash, err := hashPassword(strings.Repeat("a", 73))
	assert.ErrorIs(t, err, errPasswordTooLong)
	assert.Empty(t, hash)
}

// TestVerifyPassword は照合の判定を検証する。
func TestVerifyPassword(t *testing.T) {
	hash := mustHash(t, "correct-password")

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		// 正常系
		{name: "match", hash: hash, password: "correct-password", want: true},
		// 準正常系
		{name: "mismatch", hash: hash, password: "wrong-password", want: false},
		{name: "empty_password", hash: hash, password: "", want: false},
		{name: "case_sensitive", hash: hash, password: "Correct-Password", want: false},
		// 異常系: 保存値が平文（移行漏れ）・空・壊れたハッシュでも一致扱いにしない
		{name: "plaintext_stored", hash: "correct-password", password: "correct-password", want: false},
		{name: "empty_hash", hash: "", password: "", want: false},
		{name: "truncated_hash", hash: hash[:20], password: "correct-password", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, verifyPassword(tt.hash, tt.password))
		})
	}
}

// 異常系: ダミーハッシュは有効な bcrypt 形式で、どの入力とも一致しない（ユーザー不在時のタイミング対策用）。
func TestDummyPasswordHash(t *testing.T) {
	cost, err := bcrypt.Cost([]byte(dummyPasswordHash))
	require.NoError(t, err)
	assert.Equal(t, passwordHashCost, cost, "実ハッシュと同じコストで照合時間を揃える")

	for _, pw := range []string{"", "test-password", "password", "dummy"} {
		assert.False(t, verifyPassword(dummyPasswordHash, pw))
	}
}
