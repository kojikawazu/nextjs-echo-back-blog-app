package handlers_blogs_test

import (
	handlers_blogs "backend/handlers/blogs"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// テスト前のセットアップ
	handlers_blogs.SetupTest(nil)

	// テスト実行
	code := m.Run()

	// 終了コードを返す
	os.Exit(code)
}
