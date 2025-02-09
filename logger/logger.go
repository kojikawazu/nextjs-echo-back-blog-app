package logger

import (
	"io"
	"log"
	"os"
)

// ログレベル設定
var (
	InfoLog  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// ログ設定の初期化
func InitLogger() {
	if os.Getenv("TEST_MODE") == "true" {
		InfoLog.SetOutput(io.Discard)
		ErrorLog.SetOutput(io.Discard)
	}
}
