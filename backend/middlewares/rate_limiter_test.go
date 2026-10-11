package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

// Cloud Run のコンテナが受ける接続元（リンクローカル）を模した RemoteAddr
const testProxyRemoteAddr = "169.254.1.1:40000"

// rateLimitRequest はテストで送る 1 リクエスト分の送信元と期待ステータス。
type rateLimitRequest struct {
	xff        string // 空文字のときは X-Forwarded-For を付けない
	remoteAddr string
	wantStatus int
}

// newRateLimitTestServer は本番と同じ IP 抽出設定・エラーハンドラの Echo にミドルウェアを適用した /test を持つサーバーを返す。
func newRateLimitTestServer(m echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	e.HTTPErrorHandler = HTTPErrorHandler
	e.Use(m)
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	return e
}

func serve(e *echo.Echo, method, path string, r rateLimitRequest) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.RemoteAddr = r.remoteAddr
	if r.xff != "" {
		req.Header.Set(echo.HeaderXForwardedFor, r.xff)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// TestNewRateLimiter はクライアント IP 単位のカウントと上限超過時の応答を検証する。
// 補充を 1 時間に 1 回にして、テスト中にカウントが回復しないようにする。
func TestNewRateLimiter(t *testing.T) {
	const burst = 2

	tests := []struct {
		name     string
		requests []rateLimitRequest
	}{
		{
			// 正常系: 上限内のリクエストは通る
			name: "within_limit",
			requests: []rateLimitRequest{
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
			},
		},
		{
			// 準正常系: 上限を超えると 429 を返す
			name: "exceeded_returns_429",
			requests: []rateLimitRequest{
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusTooManyRequests},
			},
		},
		{
			// 準正常系: 異なる IP は互いのカウントに干渉しない
			name: "different_ips_are_independent",
			requests: []rateLimitRequest{
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "203.0.113.2", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "203.0.113.2", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusTooManyRequests},
			},
		},
		{
			// 準正常系: X-Forwarded-For の左端を偽装しても、右端（実際の送信元）で数えるため回避できない
			name: "spoofed_leftmost_xff_does_not_bypass",
			requests: []rateLimitRequest{
				{xff: "198.51.100.1, 203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "198.51.100.2, 203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusOK},
				{xff: "198.51.100.3, 203.0.113.1", remoteAddr: testProxyRemoteAddr, wantStatus: http.StatusTooManyRequests},
			},
		},
		{
			// 異常系: X-Forwarded-For がない場合は接続元アドレスで数える
			name: "no_xff_uses_remote_addr",
			requests: []rateLimitRequest{
				{remoteAddr: "192.0.2.1:1000", wantStatus: http.StatusOK},
				{remoteAddr: "192.0.2.1:1001", wantStatus: http.StatusOK},
				{remoteAddr: "192.0.2.2:1000", wantStatus: http.StatusOK},
				{remoteAddr: "192.0.2.1:1002", wantStatus: http.StatusTooManyRequests},
			},
		},
		{
			// 異常系: 解釈できない X-Forwarded-For は信用せず接続元アドレスで数える（値を変えても回避できない）
			name: "unparsable_xff_falls_back_to_remote_addr",
			requests: []rateLimitRequest{
				{xff: "not-an-ip-1", remoteAddr: "192.0.2.1:1000", wantStatus: http.StatusOK},
				{xff: "not-an-ip-2", remoteAddr: "192.0.2.1:1000", wantStatus: http.StatusOK},
				{xff: "not-an-ip-3", remoteAddr: "192.0.2.1:1000", wantStatus: http.StatusTooManyRequests},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := newRateLimitTestServer(newRateLimiter(rate.Every(time.Hour), burst, nil))

			for i, r := range tt.requests {
				rec := serve(e, http.MethodGet, "/test", r)
				assert.Equal(t, r.wantStatus, rec.Code, "request #%d", i)
				if r.wantStatus == http.StatusTooManyRequests {
					assert.JSONEq(t, `{"code":"TOO_MANY_REQUESTS","error":"Too many requests"}`, rec.Body.String())
				}
			}
		})
	}
}

// TestNewAPIRateLimiter は全体レートリミットがヘルスチェックを除く全パスに適用されることを検証する。
func TestNewAPIRateLimiter(t *testing.T) {
	client := rateLimitRequest{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr}

	// 正常系: ヘルスチェックは上限を超える回数でも制限されない
	t.Run("health_check_is_skipped", func(t *testing.T) {
		e := newRateLimitTestServer(newAPIRateLimiter())
		e.GET(healthCheckPath, func(c echo.Context) error {
			return c.String(http.StatusOK, "Service is running")
		})

		for i := 0; i < apiBurst*2; i++ {
			assert.Equal(t, http.StatusOK, serve(e, http.MethodGet, healthCheckPath, client).Code, "request #%d", i)
		}
	})

	// 準正常系: バースト分は通り、それを大きく超えると 429 になる
	// （補充分で数件は余分に通るため、2 倍送って 429 が含まれることを検証する）
	t.Run("exceeded_returns_429", func(t *testing.T) {
		e := newRateLimitTestServer(newAPIRateLimiter())

		for i := 0; i < apiBurst; i++ {
			assert.Equal(t, http.StatusOK, serve(e, http.MethodGet, "/test", client).Code, "request #%d", i)
		}
		limited := 0
		for i := 0; i < apiBurst; i++ {
			if serve(e, http.MethodGet, "/test", client).Code == http.StatusTooManyRequests {
				limited++
			}
		}
		assert.Positive(t, limited)
	})
}

// TestNewAuthRateLimiter は認証系レートリミットが適用先のエンドポイント間でカウントを共有し、
// 他のエンドポイントには影響しないことを検証する。
func TestNewAuthRateLimiter(t *testing.T) {
	client := rateLimitRequest{xff: "203.0.113.1", remoteAddr: testProxyRemoteAddr}

	newServer := func() *echo.Echo {
		e := echo.New()
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
		e.HTTPErrorHandler = HTTPErrorHandler
		auth := NewAuthRateLimiter()
		ok := func(c echo.Context) error { return c.String(http.StatusOK, "ok") }
		e.POST("/login", ok, auth)
		e.PUT("/update", ok, auth)
		e.GET("/blogs", ok)
		return e
	}

	// 正常系: 上限（authBurst 回）までは通る
	t.Run("within_limit", func(t *testing.T) {
		e := newServer()
		for i := 0; i < authBurst; i++ {
			assert.Equal(t, http.StatusOK, serve(e, http.MethodPost, "/login", client).Code, "request #%d", i)
		}
	})

	// 準正常系: 上限を超えたログインは 429 になる
	t.Run("exceeded_returns_429", func(t *testing.T) {
		e := newServer()
		for i := 0; i < authBurst; i++ {
			serve(e, http.MethodPost, "/login", client)
		}
		rec := serve(e, http.MethodPost, "/login", client)
		assert.Equal(t, http.StatusTooManyRequests, rec.Code)
		assert.JSONEq(t, `{"code":"TOO_MANY_REQUESTS","error":"Too many requests"}`, rec.Body.String())
	})

	// 準正常系: ログインとユーザー更新はカウントを共有する（交互に試行しても回避できない）
	t.Run("shared_between_endpoints", func(t *testing.T) {
		e := newServer()
		for i := 0; i < authBurst; i++ {
			serve(e, http.MethodPost, "/login", client)
		}
		assert.Equal(t, http.StatusTooManyRequests, serve(e, http.MethodPut, "/update", client).Code)
	})

	// 準正常系: 認証系の上限超過は、適用していないエンドポイントを止めない
	t.Run("does_not_affect_other_endpoints", func(t *testing.T) {
		e := newServer()
		for i := 0; i <= authBurst; i++ {
			serve(e, http.MethodPost, "/login", client)
		}
		assert.Equal(t, http.StatusOK, serve(e, http.MethodGet, "/blogs", client).Code)
	})
}

// TestSetupMiddlewares_IPExtractor は SetupMiddlewares がクライアント IP を X-Forwarded-For の右端から取る設定にすることを検証する。
func TestSetupMiddlewares_IPExtractor(t *testing.T) {
	tests := []struct {
		name   string
		xff    string
		wantIP string
	}{
		// 正常系: Cloud Run のフロントエンドが付与した右端の IP を使う
		{name: "single_entry", xff: "203.0.113.1", wantIP: "203.0.113.1"},
		// 準正常系: 左端に偽装した値があっても右端の IP を使う
		{name: "spoofed_leftmost", xff: "198.51.100.1, 203.0.113.1", wantIP: "203.0.113.1"},
		// 異常系: 解釈できない値が含まれる場合は接続元アドレスを使う
		{name: "unparsable", xff: "not-an-ip", wantIP: "169.254.1.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			SetupMiddlewares(e)
			e.GET("/ip", func(c echo.Context) error {
				return c.String(http.StatusOK, c.RealIP())
			})

			rec := serve(e, http.MethodGet, "/ip", rateLimitRequest{xff: tt.xff, remoteAddr: testProxyRemoteAddr})

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.wantIP, rec.Body.String())
		})
	}
}
