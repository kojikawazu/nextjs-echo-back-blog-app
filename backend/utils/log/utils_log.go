package utils

import (
	"backend/logger"
	"log"

	"github.com/labstack/echo/v4"
)

// logWithLevel はリクエスト情報（メソッド・ルートのパス・User-Agent）を添えてログを出力する。
// レベル名はロガーのプレフィックス（"INFO: " など）で付くため、メッセージには含めない。
//
// 引数:
//   - l: 出力先のロガー（logger.InfoLog など）
//   - c: Echoのリクエストコンテキスト
//   - message: 出力するログメッセージ
func logWithLevel(l *log.Logger, c echo.Context, message string) {
	l.Printf("%s %s %s - %s", c.Request().Method, c.Path(), c.Request().UserAgent(), message)
}

// LogInfo は情報ログを出力する。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - message: 出力するログメッセージ
func LogInfo(c echo.Context, message string) {
	logWithLevel(logger.InfoLog, c, message)
}

// LogError はエラーログを出力する。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - message: 出力するログメッセージ
func LogError(c echo.Context, message string) {
	logWithLevel(logger.ErrorLog, c, message)
}

// LogDebug はデバッグログを出力する。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - message: 出力するログメッセージ
func LogDebug(c echo.Context, message string) {
	logWithLevel(logger.DebugLog, c, message)
}
