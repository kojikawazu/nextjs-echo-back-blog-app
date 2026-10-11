package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// レートリミットの上限値。
// カウントはインスタンスのメモリで持つため、Cloud Run のインスタンス数だけ実効上限が緩む（--max-instances で上限を固定する）。
// BFF（フロントエンド）経由のリクエストは全員が BFF の IP に見えるため、BFF 経由の利用者全体で 1 つの枠を共有する。
const (
	// apiRate / apiBurst は全エンドポイント共通の上限（濫用によるコスト増の抑止）: 毎秒 20 件、瞬間 60 件まで
	apiRate  = rate.Limit(20)
	apiBurst = 60

	// authBurst / authInterval は認証系エンドポイントの上限（総当たり対策）: 5 回まで、以降 12 秒ごとに 1 回（= 5 回/分）
	authBurst    = 5
	authInterval = 12 * time.Second

	// rateLimiterExpiresIn はアクセスの途絶えた IP のカウントを破棄するまでの時間
	rateLimiterExpiresIn = 3 * time.Minute

	// healthCheckPath はレートリミットの対象外とするヘルスチェックのパス
	healthCheckPath = "/"
)

// errTooManyRequests は上限超過時のレスポンス。既存のエラーレスポンス形式 {"error": "..."} に合わせる。
var errTooManyRequests = &echo.HTTPError{
	Code:    http.StatusTooManyRequests,
	Message: map[string]string{"error": "Too many requests"},
}

// newAPIRateLimiter は全エンドポイント共通のレートリミットミドルウェアを返す（ヘルスチェックは対象外）。
//
// 戻り値:
//   - echo.MiddlewareFunc: レートリミットミドルウェア
func newAPIRateLimiter() echo.MiddlewareFunc {
	return newRateLimiter(apiRate, apiBurst, func(c echo.Context) bool {
		return c.Request().URL.Path == healthCheckPath
	})
}

// NewAuthRateLimiter は認証系エンドポイント（ログイン・ユーザー更新）用の厳しいレートリミットミドルウェアを返す。
// 返したインスタンスを適用したエンドポイント間でカウントを共有する。
//
// 戻り値:
//   - echo.MiddlewareFunc: レートリミットミドルウェア
func NewAuthRateLimiter() echo.MiddlewareFunc {
	return newRateLimiter(rate.Every(authInterval), authBurst, nil)
}

// newRateLimiter はクライアント IP（c.RealIP()）単位でカウントするレートリミットミドルウェアを返す。
// IP の抽出方法は Echo インスタンスの IPExtractor に従う（SetupMiddlewares で X-Forwarded-For の右端から抽出する設定にする）。
//
// 引数:
//   - r: 1 秒あたりに補充されるリクエスト数
//   - burst: 同時に許容する最大リクエスト数
//   - skipper: レートリミットを適用しない条件（nil のときは全リクエストに適用）
//
// 戻り値:
//   - echo.MiddlewareFunc: レートリミットミドルウェア
func newRateLimiter(r rate.Limit, burst int, skipper middleware.Skipper) echo.MiddlewareFunc {
	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: skipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      r,
			Burst:     burst,
			ExpiresIn: rateLimiterExpiresIn,
		}),
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return errTooManyRequests
		},
	})
}
