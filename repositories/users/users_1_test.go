package repositories_users

import (
	"backend/supabase"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupSupabase() {
	// 環境変数の読み込み
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Println("No ../../.env.test file found")
	}

	// テストの前にSupabaseクライアントの初期化
	err = supabase.InitSupabase()
	if err != nil {
		log.Fatalf("Supabase initialization failed: %v", err)
	}
}

func TestRepository_FetchUserByEmailAndPassword(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase()

	// リポジトリのインスタンスを作成
	repo := NewUserRepository()

	// テスト用の環境変数を取得
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	// メソッドを実行
	user, err := repo.FetchUserByEmailAndPassword(testEmail, testPasswd)
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
	repo := NewUserRepository()

	// メソッドを実行
	user, err := repo.FetchUserByEmailAndPassword("", "")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, user)
}
