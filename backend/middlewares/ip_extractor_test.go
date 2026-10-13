package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// newIPExtractorRequest は X-Client-IP・X-Forwarded-For・接続元を指定したリクエストを返す（空文字のヘッダは付けない）。
func newIPExtractorRequest(clientIP, xff, remoteAddr string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = remoteAddr
	if clientIP != "" {
		req.Header.Set(clientIPHeader, clientIP)
	}
	if xff != "" {
		req.Header.Set(echo.HeaderXForwardedFor, xff)
	}
	return req
}

// TestNewIPExtractor は X-Client-IP の信頼設定ごとのクライアント IP の抽出結果を検証する。
// 接続元は Cloud Run のコンテナが受ける接続元（リンクローカル）、X-Forwarded-For の右端は BFF の IP を模す。
func TestNewIPExtractor(t *testing.T) {
	const bffIP = "34.0.0.1"

	tests := []struct {
		name     string
		trust    bool
		clientIP string
		xff      string
		want     string
	}{
		{
			// 正常系: 信頼が有効なら X-Client-IP を使う
			name:     "trusted_uses_client_ip_header",
			trust:    true,
			clientIP: "203.0.113.1",
			xff:      bffIP,
			want:     "203.0.113.1",
		},
		{
			// 正常系: IPv6 は正規化した表記で返す（表記揺れで別の枠にならない）
			name:     "trusted_normalizes_ipv6",
			trust:    true,
			clientIP: "2001:DB8:0:0:0:0:0:1",
			xff:      bffIP,
			want:     "2001:db8::1",
		},
		{
			// 準正常系: 前後の空白は取り除いて解釈する
			name:     "trusted_trims_spaces",
			trust:    true,
			clientIP: " 203.0.113.1 ",
			xff:      bffIP,
			want:     "203.0.113.1",
		},
		{
			// 準正常系: 信頼が有効でもヘッダがなければ X-Forwarded-For の右端を使う
			name:  "trusted_without_header_falls_back_to_xff",
			trust: true,
			xff:   bffIP,
			want:  bffIP,
		},
		{
			// 異常系: IP として解釈できない値は使わず X-Forwarded-For の右端を使う
			name:     "trusted_unparsable_header_falls_back_to_xff",
			trust:    true,
			clientIP: "not-an-ip",
			xff:      bffIP,
			want:     bffIP,
		},
		{
			// 異常系: 複数の IP を連ねた値は使わず X-Forwarded-For の右端を使う
			name:     "trusted_multiple_ips_falls_back_to_xff",
			trust:    true,
			clientIP: "203.0.113.1, 198.51.100.1",
			xff:      bffIP,
			want:     bffIP,
		},
		{
			// 準正常系: 信頼が無効なら X-Client-IP を無視する（公開状態での偽装を防ぐ）
			name:     "untrusted_ignores_client_ip_header",
			trust:    false,
			clientIP: "203.0.113.1",
			xff:      bffIP,
			want:     bffIP,
		},
		{
			// 準正常系: 信頼が無効なら従来どおり X-Forwarded-For の左端（偽装可能）を使わない
			name:  "untrusted_ignores_spoofed_leftmost_xff",
			trust: false,
			xff:   "198.51.100.1, " + bffIP,
			want:  bffIP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newIPExtractorRequest(tt.clientIP, tt.xff, testProxyRemoteAddr)
			assert.Equal(t, tt.want, NewIPExtractor(tt.trust)(req))
		})
	}
}

// TestSetupMiddlewares_TrustClientIPHeader は環境変数 TRUST_CLIENT_IP_HEADER が "true" のときのみ X-Client-IP を信頼することを検証する。
func TestSetupMiddlewares_TrustClientIPHeader(t *testing.T) {
	const bffIP = "34.0.0.1"

	tests := []struct {
		name  string
		value string
		want  string
	}{
		// 正常系: "true" で X-Client-IP を信頼する
		{name: "true_trusts_header", value: "true", want: "203.0.113.1"},
		// 準正常系: 未設定は信頼しない（ローカル・IT/E2E の既定）
		{name: "unset_does_not_trust", value: "", want: bffIP},
		// 準正常系: "true" 以外の値は信頼しない（書き間違いは安全側に倒す）
		{name: "other_value_does_not_trust", value: "1", want: bffIP},
		{name: "uppercase_does_not_trust", value: "TRUE", want: bffIP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(trustClientIPHeaderEnv, tt.value)
			e := echo.New()
			SetupMiddlewares(e)

			req := newIPExtractorRequest("203.0.113.1", bffIP, testProxyRemoteAddr)
			assert.Equal(t, tt.want, e.IPExtractor(req))
		})
	}
}
