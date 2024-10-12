package utils

import (
	"bytes"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// LogInfo関数が正しく動作するかどうかをテスト
func TestLogInfo(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/test", nil)
	req.Header.Set("User-Agent", "TestUserAgent")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ログ出力をキャプチャするためのバッファを作成
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// テスト実行
	LogInfo(c, "This is an info message")

	// ログ出力の確認
	output := buf.String()
	// 日付と時刻を除去
	startIndex := len("2024/10/04 01:00:37 ") // 日付と時刻の長さを取得
	if len(output) > startIndex {
		outputWithoutDate := output[startIndex:] // 日付と時刻を除去した出力
		expected := "INFO: POST  TestUserAgent - This is an info message\n"
		assert.Equal(t, expected, outputWithoutDate, "ログにINFOが含まれていません")
	} else {
		t.Error("ログ出力が不正です")
	}
}

// LogError関数が正しく動作するかどうかをテスト
func TestLogError(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/test", nil)
	req.Header.Set("User-Agent", "TestUserAgent")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ログ出力をキャプチャするためのバッファを作成
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// テスト実行
	LogError(c, "This is an error message")

	// ログ出力の確認
	output := buf.String()
	// 日付と時刻を除去
	startIndex := len("2024/10/04 01:00:37 ") // 日付と時刻の長さを取得
	if len(output) > startIndex {
		outputWithoutDate := output[startIndex:] // 日付と時刻を除去した出力
		expected := "ERROR: POST  TestUserAgent - This is an error message\n"
		assert.Equal(t, expected, outputWithoutDate, "ログにINFOが含まれていません")
	} else {
		t.Error("ログ出力が不正です")
	}
}

// LogDebug関数が正しく動作するかどうかをテスト
func TestLogDebug(t *testing.T) {
	// Echoのセットアップ
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/test", nil)
	req.Header.Set("User-Agent", "TestUserAgent")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// ログ出力をキャプチャするためのバッファを作成
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// テスト実行
	LogDebug(c, "This is a debug message")

	// ログ出力の確認
	output := buf.String()
	// 日付と時刻を除去
	startIndex := len("2024/10/04 01:00:37 ") // 日付と時刻の長さを取得
	if len(output) > startIndex {
		outputWithoutDate := output[startIndex:] // 日付と時刻を除去した出力
		expected := "DEBUG: POST  TestUserAgent - This is a debug message\n"
		assert.Equal(t, expected, outputWithoutDate, "ログにINFOが含まれていません")
	} else {
		t.Error("ログ出力が不正です")
	}
}
