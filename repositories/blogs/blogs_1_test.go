package repositories_blogs

import (
	"backend/supabase"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func setupSupabase(t *testing.T) {
	// 環境変数の読み込み
	err := godotenv.Load("../../.env.test")
	if err != nil {
		t.Log("No ../../.env.test file found")
	}

	// テストの前にSupabaseクライアントの初期化
	err = supabase.InitSupabase()
	if err != nil {
		t.Fatalf("Supabase initialization failed: %v", err)
	}
}

func TestRepository_FetchBlogs(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewBlogRepository()

	// メソッドを実行
	blogs, err := repo.FetchBlogs()
	if err != nil {
		t.Fatalf("Failed to fetch blogs: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.Len(t, blogs, 1)
}

func TestRepository_FetchBlogByUserId(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewBlogRepository()

	// 環境変数からユーザIDを取得
	testUserId := os.Getenv("TEST_USER_ID")
	testBlogTitle := os.Getenv("TEST_BLOG_TITLE")

	// メソッドを実行
	blog, err := repo.FetchBlogByUserId(testUserId)
	if err != nil {
		t.Fatalf("Failed to fetch blog: %v", err)
	}

	// エラーチェックとデータ確認
	assert.NoError(t, err)
	assert.NotNil(t, blog)
	assert.Equal(t, testBlogTitle, blog.Title)
}

func TestRepository_FetchBlogByUserId_ErrorCase(t *testing.T) {
	// Supabaseクライアントの初期化
	setupSupabase(t)

	// リポジトリのインスタンスを作成
	repo := NewBlogRepository()

	// メソッドを実行
	blog, err := repo.FetchBlogByUserId("2")

	// エラーチェックとデータ確認
	assert.Error(t, err)
	assert.Nil(t, blog)
	assert.Equal(t, "ERROR: invalid input syntax for type uuid: \"2\" (SQLSTATE 22P02)", err.Error())
}
