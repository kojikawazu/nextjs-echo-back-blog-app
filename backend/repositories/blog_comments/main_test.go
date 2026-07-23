//go:build integration

package repositories_blog_comments

import (
	"backend/testsupport"
	"os"
	"testing"
)

// TestMain は testsupport 経由でテスト用 DB（testcontainers or 実 DB）を準備し、
// パッケージ内の全テストを実行する。
func TestMain(m *testing.M) {
	os.Exit(testsupport.Start(m))
}
