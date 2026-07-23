//go:build integration

// Package testsupport は Repository 層の IT（インテグレーションテスト）向けに
// PostgreSQL のテスト用 DB を提供する。
//
// 動作は環境変数 SUPABASE_URL の有無で切り替わる（判定は実環境変数のみ。
// .env ファイルの自動読み込みは行わない。ローカルの古い .env.test が
// コンテナ経路を意図せず乗っ取るのを防ぐため）:
//   - SUPABASE_URL が設定済み: その DB をそのまま使う（実 Supabase 等）。
//     この場合、呼び出し側が TEST_* も併せて指定すること。
//   - 未設定（既定）: testcontainers で使い捨ての PostgreSQL コンテナを起動し、
//     schema.sql / seed.sql を適用し、TEST_* をフィクスチャ値に設定してから接続する。
//
// build tag 'integration' を付けているため、通常の `go test ./...`（単体テスト）や
// 本番ビルドには含まれない。IT は `go test -tags=integration ./repositories/...` で実行する。
package testsupport

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"backend/logger"
	"backend/supabase"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// フィクスチャ値。testdata/seed.sql の投入データと一致させること。
const (
	fixtureUserID    = "11111111-1111-1111-1111-111111111111"
	fixtureUserName  = "Test User"
	fixtureUserEmail = "test@example.com"
	fixtureUserPass  = "test-password"
	fixtureBlogID    = "22222222-2222-2222-2222-222222222222"
)

// Start は TestMain から呼び出し、テスト用 DB を準備して m.Run() を実行する。
// 戻り値の終了コードを TestMain 側で os.Exit に渡すこと。
//
// 引数:
//   - m: テストのエントリポイント（*testing.M）
//
// 戻り値:
//   - int: m.Run() の終了コード（準備失敗時は 1）
func Start(m *testing.M) int {
	root := moduleRoot()

	logger.InitLogger()

	// SUPABASE_URL が実環境変数として指定済みなら、その DB をそのまま使う。
	if os.Getenv("SUPABASE_URL") != "" {
		if err := supabase.InitSupabase(); err != nil {
			fmt.Fprintf(os.Stderr, "testsupport: supabase init failed (SUPABASE_URL): %v\n", err)
			return 1
		}
		defer supabase.ClosePool()
		return m.Run()
	}

	// 未指定なら testcontainers で PostgreSQL を起動する。
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(
			filepath.Join(root, "testsupport", "testdata", "schema.sql"),
			filepath.Join(root, "testsupport", "testdata", "seed.sql"),
		),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "testsupport: failed to start postgres container: %v\n", err)
		return 1
	}
	defer func() { _ = container.Terminate(ctx) }()

	host, err := container.Host(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "testsupport: failed to get container host: %v\n", err)
		return 1
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "testsupport: failed to get mapped port: %v\n", err)
		return 1
	}

	// アプリの接続経路（supabase.InitSupabase）と TEST_* を環境変数で供給する。
	// コンテナは SSL 非対応なので DB_SSLMODE=disable を指定する。
	envs := map[string]string{
		"SUPABASE_URL":     fmt.Sprintf("postgresql://postgres:postgres@%s:%s/testdb", host, port.Port()),
		"DB_SSLMODE":       "disable",
		"TEST_USER_ID":     fixtureUserID,
		"TEST_USER_NAME":   fixtureUserName,
		"TEST_USER_EMAIL":  fixtureUserEmail,
		"TEST_USER_PASSWD": fixtureUserPass,
		"TEST_BLOG_ID":     fixtureBlogID,
	}
	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			fmt.Fprintf(os.Stderr, "testsupport: failed to set env %s: %v\n", k, err)
			return 1
		}
	}

	if err := supabase.InitSupabase(); err != nil {
		fmt.Fprintf(os.Stderr, "testsupport: supabase init failed (container): %v\n", err)
		return 1
	}
	defer supabase.ClosePool()

	return m.Run()
}

// moduleRoot は現在の作業ディレクトリから上方向へ go.mod を探索し、
// モジュールルート（backend/）の絶対パスを返す。
// 各テストパッケージの階層差（../../ と ../../../）に依存しないための仕組み。
//
// 戻り値:
//   - string: go.mod を含むディレクトリ。見つからない場合は開始ディレクトリ
func moduleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	start := dir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return start
		}
		dir = parent
	}
}
