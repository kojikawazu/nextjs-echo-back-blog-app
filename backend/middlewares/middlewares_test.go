package middlewares

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRequestLogger はアクセスログミドルウェアが確定したステータス・エラーを JSON で記録することを検証する。
func TestNewRequestLogger(t *testing.T) {
	tests := []struct {
		name       string
		handler    echo.HandlerFunc
		wantStatus int
		wantLevel  string
		wantError  string // 空文字のときは error キーが存在しないことを検証する
	}{
		{
			// 正常系: 成功レスポンスは INFO で記録される
			name: "success_logs_info",
			handler: func(c echo.Context) error {
				return c.String(http.StatusOK, "ok")
			},
			wantStatus: http.StatusOK,
			wantLevel:  "INFO",
		},
		{
			// 準正常系: echo.HTTPError（404）はエラーハンドラ適用後のステータスで記録される
			name: "http_error_logs_resolved_status",
			handler: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusNotFound, "not found")
			},
			wantStatus: http.StatusNotFound,
			wantLevel:  "ERROR",
			wantError:  "code=404, message=not found",
		},
		{
			// 準正常系: ハンドラが 400 を明示的に返した場合はエラーなしで記録される
			name: "bad_request_response_logs_info",
			handler: func(c echo.Context) error {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid"})
			},
			wantStatus: http.StatusBadRequest,
			wantLevel:  "INFO",
		},
		{
			// 異常系: 想定外のエラーは 500 として ERROR で記録される
			name: "unexpected_error_logs_500",
			handler: func(c echo.Context) error {
				return errors.New("unexpected failure")
			},
			wantStatus: http.StatusInternalServerError,
			wantLevel:  "ERROR",
			wantError:  "unexpected failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			e := echo.New()
			e.Use(newRequestLogger(&buf))
			e.GET("/test", tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/test?q=1", nil)
			req.Header.Set("User-Agent", "test-agent")
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			var entry map[string]any
			require.NoError(t, json.Unmarshal(buf.Bytes(), &entry), "log must be a single JSON line: %q", buf.String())

			assert.Equal(t, tt.wantLevel, entry["level"])
			assert.Equal(t, "request", entry["msg"])
			assert.Equal(t, http.MethodGet, entry["method"])
			assert.Equal(t, "/test?q=1", entry["uri"])
			assert.Equal(t, "test-agent", entry["user_agent"])
			assert.Equal(t, float64(tt.wantStatus), entry["status"])

			if tt.wantError == "" {
				assert.NotContains(t, entry, "error")
			} else {
				assert.Equal(t, tt.wantError, entry["error"])
			}
		})
	}
}

// TestNewRequestLogger_DoesNotLogSensitiveHeaders は Authorization ヘッダ等がログに含まれないことを検証する（準正常系）。
func TestNewRequestLogger_DoesNotLogSensitiveHeaders(t *testing.T) {
	var buf bytes.Buffer
	e := echo.New()
	e.Use(newRequestLogger(&buf))
	e.GET("/test", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer secret-token-value")
	req.Header.Set("Cookie", "token=secret-cookie-value")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.NotContains(t, buf.String(), "secret-token-value")
	assert.NotContains(t, buf.String(), "secret-cookie-value")
}
