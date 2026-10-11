package routes

import (
	"backend/middlewares"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// newRoutesTestServer は本番と同じ IP 抽出設定・エラーハンドラでルーティングを設定した Echo を返す。
// レートリミットはハンドラより先に評価されるため、DB に到達しないリクエスト（不正ボディ・未認証）で検証できる。
func newRoutesTestServer() *echo.Echo {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.HTTPErrorHandler = middlewares.HTTPErrorHandler
	SetupRoutes(e)
	return e
}

func serveRoute(e *echo.Echo, method, path, xff string) int {
	req := httptest.NewRequest(method, path, strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderXForwardedFor, xff)
	req.RemoteAddr = "169.254.1.1:40000"
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code
}

// TestSetupRoutes_AuthRateLimit は認証系エンドポイントに厳しいレートリミットが適用されていることを検証する。
func TestSetupRoutes_AuthRateLimit(t *testing.T) {
	const client = "203.0.113.1"
	// 認証系の上限（middlewares の authBurst）
	const authLimit = 5

	// 正常系: 上限内のログイン試行はハンドラまで届く（不正ボディのため 400）
	t.Run("login_within_limit_reaches_handler", func(t *testing.T) {
		e := newRoutesTestServer()
		for i := 0; i < authLimit; i++ {
			assert.Equal(t, http.StatusBadRequest, serveRoute(e, http.MethodPost, "/api/users/login", client), "request #%d", i)
		}
	})

	// 準正常系: 上限を超えたログイン試行は 429 になる
	t.Run("login_exceeded_returns_429", func(t *testing.T) {
		e := newRoutesTestServer()
		for i := 0; i < authLimit; i++ {
			serveRoute(e, http.MethodPost, "/api/users/login", client)
		}
		assert.Equal(t, http.StatusTooManyRequests, serveRoute(e, http.MethodPost, "/api/users/login", client))
	})

	// 準正常系: ユーザー更新はログインとカウントを共有する
	t.Run("update_shares_limit_with_login", func(t *testing.T) {
		e := newRoutesTestServer()
		for i := 0; i < authLimit; i++ {
			serveRoute(e, http.MethodPost, "/api/users/login", client)
		}
		assert.Equal(t, http.StatusTooManyRequests, serveRoute(e, http.MethodPut, "/api/users/update", client))
	})

	// 準正常系: 認証系の上限超過後も、制限対象外のエンドポイント（ログアウト）は通常どおり処理される
	t.Run("other_endpoints_unaffected", func(t *testing.T) {
		e := newRoutesTestServer()
		for i := 0; i <= authLimit; i++ {
			serveRoute(e, http.MethodPost, "/api/users/login", client)
		}
		assert.Equal(t, http.StatusOK, serveRoute(e, http.MethodPost, "/api/users/logout", client))
	})
}
