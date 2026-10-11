//go:build e2e

package e2e

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 認証系レートリミットの上限（middlewares の authBurst）
const authRateLimit = 5

// 準正常系: 同じ送信元からのログイン試行は上限を超えると 429 になる。
func TestE2E_Login_RateLimited(t *testing.T) {
	const clientIP = "203.0.113.10"
	client := newClient(t)
	wrongLogin := map[string]string{
		"email":    os.Getenv("TEST_USER_EMAIL"),
		"password": "wrong-password",
	}

	for i := 0; i < authRateLimit; i++ {
		resp := doJSONFrom(t, client, clientIP, http.MethodPost, "/api/users/login", wrongLogin)
		resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "attempt #%d", i)
	}

	var body map[string]string
	resp := doJSONFrom(t, client, clientIP, http.MethodPost, "/api/users/login", wrongLogin)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	decodeBody(t, resp, &body)
	assert.Equal(t, "Too many requests", body["error"])
}

// 準正常系: 上限に達した送信元でも、正しいパスワードでのログインは通らない（総当たりの途中で正解しても弾く）。
// 別の送信元からのログインは影響を受けない。
func TestE2E_Login_RateLimitIsPerClient(t *testing.T) {
	const blockedIP = "203.0.113.11"
	client := newClient(t)

	for i := 0; i < authRateLimit; i++ {
		resp := doJSONFrom(t, client, blockedIP, http.MethodPost, "/api/users/login", map[string]string{
			"email":    os.Getenv("TEST_USER_EMAIL"),
			"password": "wrong-password",
		})
		resp.Body.Close()
	}

	validLogin := map[string]string{
		"email":    os.Getenv("TEST_USER_EMAIL"),
		"password": os.Getenv("TEST_USER_PASSWD"),
	}

	blocked := doJSONFrom(t, client, blockedIP, http.MethodPost, "/api/users/login", validLogin)
	blocked.Body.Close()
	assert.Equal(t, http.StatusTooManyRequests, blocked.StatusCode)

	other := doJSONFrom(t, client, "203.0.113.12", http.MethodPost, "/api/users/login", validLogin)
	other.Body.Close()
	assert.Equal(t, http.StatusOK, other.StatusCode)
}
