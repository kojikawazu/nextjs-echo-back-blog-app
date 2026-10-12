package utils

import (
	"backend/logger"
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// captureLoggers は logger の各ロガーの出力をバッファに差し替え、日時・ファイル名のフラグを外す。
// テスト終了時に元の出力先・フラグへ戻す。
func captureLoggers(t *testing.T) map[*log.Logger]*bytes.Buffer {
	t.Helper()
	bufs := map[*log.Logger]*bytes.Buffer{}
	for _, l := range []*log.Logger{logger.InfoLog, logger.ErrorLog, logger.WarnLog, logger.DebugLog, logger.TestLog} {
		buf := &bytes.Buffer{}
		writer, flags := l.Writer(), l.Flags()
		l.SetOutput(buf)
		l.SetFlags(0)
		t.Cleanup(func() {
			l.SetOutput(writer)
			l.SetFlags(flags)
		})
		bufs[l] = buf
	}
	return bufs
}

// newContext はメソッド・ルートのパス・User-Agent を指定した Echo のコンテキストを返す。
func newContext(path, userAgent string) echo.Context {
	req := httptest.NewRequest(http.MethodPost, "/api/blogs/1", nil)
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	c := echo.New().NewContext(req, httptest.NewRecorder())
	c.SetPath(path)
	return c
}

// TestLog は各レベルのログが対応するロガーにだけ、レベル名を 1 回だけ付けて出力されることを検証する。
func TestLog(t *testing.T) {
	tests := []struct {
		name       string
		logFunc    func(echo.Context, string)
		path       string
		userAgent  string
		message    string
		wantLogger *log.Logger
		want       string
	}{
		{
			// 正常系: INFO はプレフィックスだけにレベル名が付き、ルートのパターンが出る
			name:       "info",
			logFunc:    LogInfo,
			path:       "/api/blogs/:id",
			userAgent:  "TestUserAgent",
			message:    "fetch blog",
			wantLogger: logger.InfoLog,
			want:       "INFO: POST /api/blogs/:id TestUserAgent - fetch blog\n",
		},
		{
			// 正常系: ERROR は ErrorLog に出る
			name:       "error",
			logFunc:    LogError,
			path:       "/api/blogs/:id",
			userAgent:  "TestUserAgent",
			message:    "failed",
			wantLogger: logger.ErrorLog,
			want:       "ERROR: POST /api/blogs/:id TestUserAgent - failed\n",
		},
		{
			// 正常系: DEBUG は DebugLog に出る
			name:       "debug",
			logFunc:    LogDebug,
			path:       "/api/blogs/:id",
			userAgent:  "TestUserAgent",
			message:    "detail",
			wantLogger: logger.DebugLog,
			want:       "DEBUG: POST /api/blogs/:id TestUserAgent - detail\n",
		},
		{
			// 準正常系: User-Agent が無い
			name:       "no_user_agent",
			logFunc:    LogInfo,
			path:       "/api/blogs/:id",
			userAgent:  "",
			message:    "fetch blog",
			wantLogger: logger.InfoLog,
			want:       "INFO: POST /api/blogs/:id  - fetch blog\n",
		},
		{
			// 準正常系: ルートに一致せずパスが空
			name:       "no_route",
			logFunc:    LogInfo,
			path:       "",
			userAgent:  "TestUserAgent",
			message:    "fetch blog",
			wantLogger: logger.InfoLog,
			want:       "INFO: POST  TestUserAgent - fetch blog\n",
		},
		{
			// 準正常系: メッセージが空
			name:       "empty_message",
			logFunc:    LogInfo,
			path:       "/api/blogs/:id",
			userAgent:  "TestUserAgent",
			message:    "",
			wantLogger: logger.InfoLog,
			want:       "INFO: POST /api/blogs/:id TestUserAgent - \n",
		},
		{
			// 準正常系: メッセージ中の % は書式指定として解釈されない
			name:       "percent_in_message",
			logFunc:    LogInfo,
			path:       "/api/blogs/:id",
			userAgent:  "TestUserAgent",
			message:    "100% done %s %d",
			wantLogger: logger.InfoLog,
			want:       "INFO: POST /api/blogs/:id TestUserAgent - 100% done %s %d\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufs := captureLoggers(t)

			tt.logFunc(newContext(tt.path, tt.userAgent), tt.message)

			for l, buf := range bufs {
				if l == tt.wantLogger {
					assert.Equal(t, tt.want, buf.String())
					continue
				}
				// 準正常系: 他のレベルのロガーには出力しない
				assert.Empty(t, buf.String(), "prefix=%q", l.Prefix())
			}
		})
	}
}
