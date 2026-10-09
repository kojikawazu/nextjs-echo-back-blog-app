package middlewares

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupMiddlewares は Echo にロガー・リカバリー・CORS の各ミドルウェアを設定する。
//
// 引数:
//   - e: ミドルウェアを設定する対象の Echo インスタンス
func SetupMiddlewares(e *echo.Echo) {
	// ロガーとリカバリーミドルウェアを使用
	e.Use(newRequestLogger(os.Stdout))
	e.Use(middleware.Recover())

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// CORSを有効化
	// AllowCredentialsをtrueに設定すると、クライアント側でwithCredentialsをtrueに設定する必要がある
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(allowedOrigins, ","),
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAuthorization,
			echo.HeaderAccessControlAllowCredentials,
		},
		// ExposeHeaders: []string{
		// 	echo.HeaderSetCookie,
		// },
		AllowCredentials: true,
	}))
}

// newRequestLogger は HTTP リクエスト/レスポンスのアクセスログを JSON 形式で出力するミドルウェアを返す。
// 非推奨となった middleware.Logger の後継として middleware.RequestLoggerWithConfig を使用する。
// リクエストヘッダ・ボディ等のセンシティブ情報は記録しない。
//
// 引数:
//   - w: ログの出力先
//
// 戻り値:
//   - echo.MiddlewareFunc: アクセスログミドルウェア
func newRequestLogger(w io.Writer) echo.MiddlewareFunc {
	logger := slog.New(slog.NewJSONHandler(w, nil))

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// ハンドラのエラーをエラーハンドラに渡し、確定したステータスを記録する
		HandleError:      true,
		LogLatency:       true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("remote_ip", v.RemoteIP),
				slog.String("host", v.Host),
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.String("user_agent", v.UserAgent),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("bytes_in", v.ContentLength),
				slog.Int64("bytes_out", v.ResponseSize),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("error", v.Error.Error()))
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "request", attrs...)
				return nil
			}
			logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "request", attrs...)
			return nil
		},
	})
}
