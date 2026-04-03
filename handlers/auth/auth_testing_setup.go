package handlers_auth

import (
	"backend/logger"
	"testing"

	"github.com/joho/godotenv"
)

// SetupTest はテストの前に環境変数を読み込み、ログ設定を初期化します
func SetupTest(t *testing.T) {
	err := godotenv.Load("../../.env.test")
	if err != nil && t != nil {
		t.Log("No ../../.env.test file found")
	}

	logger.InitLogger()
}
