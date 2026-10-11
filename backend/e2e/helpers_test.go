//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync/atomic"
	"testing"
)

// clientIPSeq は doJSON がリクエストごとに割り当てる送信元 IP の連番。
var clientIPSeq atomic.Uint32

// nextClientIP はリクエストごとに異なる送信元 IP（198.18.0.0/15 のベンチマーク用アドレス）を返す。
// E2E は全リクエストが 127.0.0.1 から届くため、そのままではレートリミットのカウントを全テストで共有してしまう。
// レートリミットの検証以外のテストがカウントの影響を受けないよう、リクエストごとに別の IP として扱わせる。
func nextClientIP() string {
	n := clientIPSeq.Add(1)
	return fmt.Sprintf("198.18.%d.%d", (n>>8)&0xff, n&0xff)
}

// newClient は Cookie を保持する HTTP クライアントを生成する。
// ログインで受け取る token / visit-id-token Cookie を後続リクエストへ引き継ぐために用いる。
func newClient(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("failed to create cookie jar: %v", err)
	}
	return &http.Client{Jar: jar}
}

// doJSON は baseURL に対して method/path のリクエストを送る。
// body が非 nil の場合は JSON エンコードして送信する。レスポンスは呼び出し側で Close する。
// 送信元 IP はリクエストごとに変える（nextClientIP 参照）。
func doJSON(t *testing.T, client *http.Client, method, path string, body any) *http.Response {
	t.Helper()
	return doJSONFrom(t, client, nextClientIP(), method, path, body)
}

// doJSONFrom は送信元 IP（X-Forwarded-For）を clientIP に固定して doJSON と同じリクエストを送る。
func doJSONFrom(t *testing.T, client *http.Client, clientIP, method, path string, body any) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, reader)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-Forwarded-For", clientIP)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed (%s %s): %v", method, path, err)
	}
	return resp
}

// decodeBody はレスポンスボディを dst にデコードし、ボディを Close する。
func decodeBody(t *testing.T, resp *http.Response, dst any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

// login はシード済みユーザーでログインし、token Cookie を client に保存する。
func login(t *testing.T, client *http.Client) {
	t.Helper()
	resp := doJSON(t, client, http.MethodPost, "/api/users/login", map[string]string{
		"email":    os.Getenv("TEST_USER_EMAIL"),
		"password": os.Getenv("TEST_USER_PASSWD"),
	})
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login failed: status=%d", resp.StatusCode)
	}
}
