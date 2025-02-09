package repositories_blogs_test

import (
	repositories_blogs "backend/repositories/blogs"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// テスト前のセットアップ
	repositories_blogs.SetupTest(nil)

	// テスト実行
	code := m.Run()

	// 終了コードを返す
	os.Exit(code)
}
