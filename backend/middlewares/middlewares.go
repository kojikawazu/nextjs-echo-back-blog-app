package middlewares

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ミドルウェアの設定
func SetupMiddlewares(e *echo.Echo) {
	// ロガーとリカバリーミドルウェアを使用
	e.Use(middleware.Logger())
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
