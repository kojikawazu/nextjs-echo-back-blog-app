package utils

import (
	"log"

	"github.com/labstack/echo/v4"
)

// リクエスト情報を共通して取得
func logWithLevel(c echo.Context, level string, message string) {
	requestPath := c.Path()
	method := c.Request().Method
	userAgent := c.Request().UserAgent()

	// レベルに応じたログを出力
	log.Printf("%s: %s %s %s - %s", level, method, requestPath, userAgent, message)
}

// 情報ログ
func LogInfo(c echo.Context, message string) {
	logWithLevel(c, "INFO", message)
}

// エラーログ
func LogError(c echo.Context, message string) {
	logWithLevel(c, "ERROR", message)
}

// デバッグログ
func LogDebug(c echo.Context, message string) {
	logWithLevel(c, "DEBUG", message)
}
