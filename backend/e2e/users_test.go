//go:build e2e

package e2e

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readBody はレスポンスボディを文字列で読み出し、ボディを Close する。
func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(b)
}

// 正常系: ログイン後のユーザー取得は公開フィールドのみを返し、password キーを含まない。
func TestE2E_UserDetail_ExcludesPassword(t *testing.T) {
	client := newClient(t)
	login(t, client)

	resp := doJSON(t, client, http.MethodGet, "/api/users/detail", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	decodeBody(t, resp, &body)
	keys := make([]string, 0, len(body))
	for k := range body {
		keys = append(keys, k)
	}
	assert.ElementsMatch(t, []string{"id", "name", "email", "created_at", "updated_at"}, keys)
	assert.Equal(t, os.Getenv("TEST_USER_ID"), body["id"])
	assert.Equal(t, os.Getenv("TEST_USER_EMAIL"), body["email"])
}

// 準正常系: 未認証のユーザー取得は 401 で、ユーザー情報を返さない。
func TestE2E_UserDetail_Unauthenticated(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodGet, "/api/users/detail", nil)
	body := readBody(t, resp)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotContains(t, body, os.Getenv("TEST_USER_EMAIL"))
}
