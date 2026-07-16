package repositories_blog_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchUserByEmailAndPassword(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewBlogUsersRepository()

	// テスト用の環境変数を取得
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	// メソッドを実行
	user, err := repo.FetchBlogUsersByEmailAndPassword(testEmail, testPasswd)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

func TestRepository_FetchUserByEmailAndPassword_ErrorCases(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewBlogUsersRepository()

	// メソッドを実行
	user, err := repo.FetchBlogUsersByEmailAndPassword("", "")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, user)
}
