//go:build e2e

package e2e

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: ヘルスチェック GET / が 200 で "Service is running" を返す。
func TestE2E_Health(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodGet, "/", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Service is running", string(body))
}
