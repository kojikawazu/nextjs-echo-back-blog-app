//go:build e2e

package e2e

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"backend/middlewares"
	"backend/routes"
	"backend/testsupport"

	"github.com/labstack/echo/v4"
)

// baseURL は httptest で起動した E2E 用サーバーのベース URL。各テストから参照する。
var baseURL string

// TestMain は E2E の前準備を行う:
//  1. testsupport.Setup でテスト用 DB（既定は testcontainers）を準備し supabase.Pool を接続
//  2. middlewares/routes で本番同等の Echo アプリを組み立て
//  3. httptest でサーバーを起動し baseURL を公開
//
// 注意: config パッケージが起動時（init）に JWT_SECRET_KEY を要求するため、
// バイナリ起動時に環境変数として JWT_SECRET_KEY を渡すこと。
func TestMain(m *testing.M) {
	teardown, err := testsupport.Setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, "e2e: DB setup failed: %v\n", err)
		os.Exit(1)
	}

	// 本番と同じ組み立て（ミドルウェア + ルーティング）。Repository は supabase.Pool を参照する。
	e := echo.New()
	middlewares.SetupMiddlewares(e)
	routes.SetupRoutes(e)

	srv := httptest.NewServer(e)
	baseURL = srv.URL

	code := m.Run()

	// os.Exit は defer を実行しないため明示的に後始末する。
	srv.Close()
	teardown()
	os.Exit(code)
}
