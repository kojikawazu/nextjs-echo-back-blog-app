package repositories_blogs

import (
	"backend/supabase"
	"testing"

	"github.com/joho/godotenv"
)

// setupSupabase はテストの前にSupabaseクライアントを初期化します
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
