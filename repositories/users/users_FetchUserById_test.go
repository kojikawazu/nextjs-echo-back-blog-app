package repositories_users

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FetchUserById(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// テスト用の環境変数を取得
	testUserId := os.Getenv("TEST_USER_ID")
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")

	// メソッドを実行
	user, err := repo.FetchUserById(testUserId)
	if err != nil {
		t.Fatalf("Failed to fetch user by id: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

func TestRepository_FetchUserById_InvalidId(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// メソッドを実行
	user, err := repo.FetchUserById("")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, user)
}
