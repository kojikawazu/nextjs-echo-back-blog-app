package handlers_auth

import (
	"backend/logger"
	"testing"

	"github.com/joho/godotenv"
)

// SetupTest はテストの前に環境変数を読み込み、ログ設定を初期化します
//
// 引数:
//   - t: 環境変数ファイルが無い場合のログ出力に用いるテスト構造体
func SetupTest(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	if err != nil && t != nil {
		t.Log("No ../../.env.test file found")
	}

	logger.InitLogger()
}
