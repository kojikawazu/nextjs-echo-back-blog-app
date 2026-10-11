//go:build integration

package repositories_blogs_test

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

// notExistBlogID はシードに存在しない UUID 形式のブログ ID（不在時の挙動の検証に用いる）。
const notExistBlogID = "00000000-0000-0000-0000-000000000000"
