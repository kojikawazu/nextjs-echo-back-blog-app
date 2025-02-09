package services_blogs_test

import (
	services_blogs "backend/services/blogs"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// テスト前のセットアップ
	services_blogs.SetupTest(nil)

	// テスト実行
	code := m.Run()

	// 終了コードを返す
	os.Exit(code)
}
