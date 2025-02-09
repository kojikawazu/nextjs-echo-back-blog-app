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
	WarnLog  = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	TestLog  = log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// ログ設定の初期化
func InitLogger() {
	if os.Getenv("TEST_MODE") == "true" {
		InfoLog.SetOutput(io.Discard)
		ErrorLog.SetOutput(io.Discard)
		WarnLog.SetOutput(io.Discard)
		DebugLog.SetOutput(io.Discard)
	}
}
