package repositories_blogs

import (
	"backend/logger"
	"backend/supabase"
	"testing"

	"github.com/joho/godotenv"
)

// SetupTest はテストの前に環境変数を読み込み、ログ設定を初期化します
func SetupTest(t *testing.T) {
	// 環境変数の読み込み
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		t.Log("No ../../../.env.test file found")
	}

	// ログ設定の初期化
	logger.InitLogger()

	// テストの前にSupabaseクライアントの初期化
	err = supabase.InitSupabase()
	if err != nil {
		t.Fatalf("Supabase initialization failed: %v", err)
	}

}
