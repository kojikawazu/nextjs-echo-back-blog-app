//go:build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
)

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
func doJSON(t *testing.T, client *http.Client, method, path string, body any) *http.Response {
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
