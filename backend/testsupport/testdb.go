//go:build integration || e2e

// Package testsupport は Repository 層の IT および API-E2E 向けに
// PostgreSQL のテスト用 DB を提供する。
//
// 動作は環境変数 SUPABASE_URL の有無で切り替わる（判定は実環境変数のみ。
// .env ファイルの自動読み込みは行わない。ローカルの古い .env.test が
// コンテナ経路を意図せず乗っ取るのを防ぐため）:
//   - SUPABASE_URL が設定済み: その DB をそのまま使う（実 Supabase 等）。
//     この場合、呼び出し側が TEST_* も併せて指定すること。
//   - 未設定（既定）: testcontainers で使い捨ての PostgreSQL を起動し、
//     schema.sql / seed.sql を適用し、TEST_* をフィクスチャ値に設定してから接続する。
//
// build tag 'integration'（IT）/ 'e2e'（E2E）を付けているため、通常の
// `go test ./...`（単体テスト）や本番ビルドには含まれない。
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
// 戻り値の終了コードを TestMain 側で os.Exit に渡すこと（IT 向けの薄いラッパ）。
//
// 引数:
//   - m: テストのエントリポイント（*testing.M）
//
// 戻り値:
//   - int: m.Run() の終了コード（準備失敗時は 1）
func Start(m *testing.M) int {
	teardown, err := Setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "testsupport: %v\n", err)
		return 1
	}
	defer teardown()
	return m.Run()
}

// Setup はテスト用 DB を準備し、後始末を行う teardown 関数を返す。
// E2E のように m.Run() の前後で追加のセットアップ（サーバ起動等）を挟みたい
// 場合に用いる。DB 準備に失敗した場合はエラーを返す（コンテナは破棄済み）。
//
// 戻り値:
//   - func(): プール close とコンテナ破棄を行う後始末関数（成功時のみ非 nil）
//   - error: DB 準備に失敗した場合のエラー
func Setup() (func(), error) {
	root := moduleRoot()

	logger.InitLogger()

	// SUPABASE_URL が実環境変数として指定済みなら、その DB をそのまま使う。
	if os.Getenv("SUPABASE_URL") != "" {
		if err := supabase.InitSupabase(); err != nil {
			return nil, fmt.Errorf("supabase init failed (SUPABASE_URL): %w", err)
		}
		return func() { supabase.ClosePool() }, nil
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
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}
	terminate := func() { _ = container.Terminate(ctx) }

	host, err := container.Host(ctx)
	if err != nil {
		terminate()
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		terminate()
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
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
			terminate()
			return nil, fmt.Errorf("failed to set env %s: %w", k, err)
		}
	}

	if err := supabase.InitSupabase(); err != nil {
		terminate()
		return nil, fmt.Errorf("supabase init failed (container): %w", err)
	}

	return func() {
		supabase.ClosePool()
		terminate()
	}, nil
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
