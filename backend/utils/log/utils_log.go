package utils

import (
	"backend/logger"
	"log"

	"github.com/labstack/echo/v4"
)

// リクエスト情報を共通して取得
func logWithLevel(c echo.Context, level string, message string) {
	requestPath := c.Path()
	method := c.Request().Method
	userAgent := c.Request().UserAgent()

	// レベルに応じたログを出力
	switch level {
	case "INFO":
		logger.InfoLog.Printf("%s: %s %s %s - %s", level, method, requestPath, userAgent, message)
	case "ERROR":
		logger.ErrorLog.Printf("%s: %s %s %s - %s", level, method, requestPath, userAgent, message)
	case "WARN":
		logger.WarnLog.Printf("%s: %s %s %s - %s", level, method, requestPath, userAgent, message)
	case "DEBUG":
		logger.DebugLog.Printf("%s: %s %s %s - %s", level, method, requestPath, userAgent, message)
	case "TEST":
		logger.TestLog.Printf("%s: %s %s %s - %s", level, method, requestPath, userAgent, message)
	default:
		log.Printf("%s: %s %s %s - %s", level, method, requestPath, userAgent, message)
	}
}

// LogInfo は情報ログを出力する。
func LogInfo(c echo.Context, message string) {
	logWithLevel(c, "INFO", message)
}

// LogError はエラーログを出力する。
func LogError(c echo.Context, message string) {
	logWithLevel(c, "ERROR", message)
}

// LogDebug はデバッグログを出力する。
func LogDebug(c echo.Context, message string) {
	logWithLevel(c, "DEBUG", message)
}
