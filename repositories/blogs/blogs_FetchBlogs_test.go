package repositories_blogs

import (
	"backend/supabase"
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
	assert.Len(t, blogs, 2)
}
