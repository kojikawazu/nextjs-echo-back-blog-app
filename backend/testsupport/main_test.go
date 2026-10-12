//go:build integration || e2e

package testsupport

import (
	"os"
	"testing"
)

// TestMain はテスト用 DB を準備し、パッケージ内の全テストを実行する。
// スキーマの検査（schema_test.go）が testdata/schema.sql を適用した DB を必要とするため。
func TestMain(m *testing.M) {
	os.Exit(Start(m))
}
